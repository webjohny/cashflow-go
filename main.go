package main

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/config"
	"github.com/webjohny/cashflow-go/controller"
	"github.com/webjohny/cashflow-go/middleware"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/service"
	"gorm.io/gorm"
	"os"
)

var (
	db *gorm.DB = config.SetupDatabaseConnection()

	// Repositories
	raceRepository       repository.RaceRepository        = repository.NewRaceRepository(db)
	lobbyRepository      repository.LobbyRepository       = repository.NewLobbyRepository(db)
	userRepository       repository.UserRepository        = repository.NewUserRepository(db)
	playerRepository     repository.PlayerRepository      = repository.NewPlayerRepository(db)
	professionRepository repository.ProfessionRepository  = repository.NewProfessionRepository(os.Getenv("PROFESSIONS_PATH"))
	trxRepository        repository.TransactionRepository = repository.NewTransactionRepository(db)

	// Services
	jwtService         service.JWTService         = service.NewJWTService()
	userService        service.UserService        = service.NewUserService(userRepository)
	transactionService service.TransactionService = service.NewTransactionService(trxRepository)
	playerService      service.PlayerService      = service.NewPlayerService(playerRepository, professionRepository, transactionService)
	authService        service.AuthService        = service.NewAuthService(userRepository)
	gameService        service.GameService        = service.NewGameService(lobbyRepository, raceRepository, playerRepository)
	raceService        service.RaceService        = service.NewRaceService(raceRepository, playerService, transactionService)
	lobbyService       service.LobbyService       = service.NewLobbyService(lobbyRepository)
	cardService        service.CardService        = service.NewCardService(gameService, raceService)
	financeService     service.FinanceService     = service.NewFinanceService(raceService, playerService)

	// Controllers
	gameController    controller.GameController    = controller.NewGameController(gameService)
	lobbyController   controller.LobbyController   = controller.NewLobbyController(lobbyService)
	financeController controller.FinanceController = controller.NewFinanceController(financeService)
	cardController    controller.CardController    = controller.NewCardController(cardService)
	authController    controller.AuthController    = controller.NewAuthController(authService, jwtService)
	userController    controller.UserController    = controller.NewUserController(userService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	cardRoutes := r.Group("api/card")
	{
		cardRoutes.GET("/prepare/:family/:type", cardController.Prepare)
		cardRoutes.GET("/skip/:family/:type", cardController.Skip)
		cardRoutes.GET("/sell/:family/:type", cardController.Selling)
		cardRoutes.GET("/buy/:family/:type", cardController.Purchase)
		cardRoutes.GET("/ok/:family/:type", cardController.Accept)
	}

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
		// userRoutes.POST("/picture", userController.SaveFile)
	}

	lobbyRoutes := r.Group("api/lobby", middleware.AuthorizeJWT(jwtService))
	{
		lobbyRoutes.POST("/create", lobbyController.CreateLobby)
		lobbyRoutes.GET("/join/:lobbyId", lobbyController.Join)
		lobbyRoutes.GET("/leave/:lobbyId", lobbyController.Leave)
		// userRoutes.POST("/picture", userController.SaveFile)
	}

	gameRoutes := r.Group("api/game", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		gameRoutes.GET("", gameController.GetGame)
		gameRoutes.POST("/start", gameController.Start)
		// userRoutes.POST("/picture", userController.SaveFile)
	}

	financeRoutes := r.Group("api/finance", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		financeRoutes.POST("/send/money", financeController.SendMoney)
		financeRoutes.POST("/send/assets", financeController.SendAssets)
		financeRoutes.POST("/loan/take", financeController.TakeLoan)
	}

	// cdnRoutes := r.Group("api/cdn")
	// {
	// 	cdnRoutes.GET("/picture/:file_name", userController.GetFile)
	// }
	// userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	// {
	// 	userRoutes.GET("/profile", userController.Profile)
	// 	userRoutes.PUT("/profile", userController.Update)
	// }

	// noteRoutes := r.Group("api/notes", middleware.AuthorizeJWT(jwtService))
	// {
	// 	noteRoutes.GET("/", noteController.All)
	// 	noteRoutes.POST("/", noteController.Insert)
	// 	noteRoutes.GET("/:id", noteController.FindById)
	// 	noteRoutes.PUT("/:id", noteController.Update)
	// 	noteRoutes.DELETE("/:id", noteController.Delete)
	// }

	// pagerRoutes := r.Group("pager", middleware.AuthorizeJWT(jwtService))
	// {
	// 	pagerRoutes.GET("/", pagerController.All)
	// 	pagerRoutes.POST("/", pagerController.Insert)
	// 	pagerRoutes.GET("/:id", pagerController.FindById)
	// 	pagerRoutes.PUT("/:id", pagerController.Update)
	// 	pagerRoutes.DELETE("/:id", pagerController.Delete)
	// }

	// pagerStatusRoutes := r.Group("/api/status")
	// {
	// 	pagerStatusRoutes.GET("/:id", pagerController.FindStatusById)
	// }

	r.Run(":" + os.Getenv("PORT"))
}
