package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/config"
	"github.com/webjohny/cashflow-go/controller"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/middleware"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/storage"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
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
	gameService        service.GameService        = service.NewGameService(raceService, playerService, lobbyRepository, professionRepository)
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

var globalToken string

func switchUser(ctx *gin.Context) {
	var token *jwt.Token

	if globalToken != "" {
		token, _ = jwtService.ValidateToken(globalToken)
	}

	var email string
	password := "qwerty"

	if token != nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)

		if claims["user_id"] == "1" {
			email = "webtoolteam@gmail.com"
		}
	}

	if email == "" {
		email = "geryh213921@gmail.com"
	}

	authResult := authService.VerifyCredential(email, password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := jwtService.GenerateToken(strconv.FormatUint(v.ID, 10), v.Email, v.Profile, v.Jk, v.Name)
		v.Token = generatedToken

		globalToken = v.Token
		response := request.BuildResponse(true, "OK", v)
		ctx.JSON(http.StatusOK, response)
		return
	} else {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Errorf(storage.ErrorForbidden).Error(),
		})
	}
}

func resetToken() {

}

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	cardRoutes := r.Group("api/card", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		cardRoutes.GET("/type", cardController.Type)
		cardRoutes.GET("/reset-transaction", cardController.ResetTransaction)
		cardRoutes.POST("/prepare/:family/:type", cardController.Prepare)
		cardRoutes.POST("/skip/:family/:type", cardController.Skip)
		cardRoutes.POST("/sell/:family/:type", cardController.Selling)
		cardRoutes.POST("/buy/:family/:type", cardController.Purchase)
		cardRoutes.POST("/ok/:family/:type", cardController.Accept)
	}

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/switch", switchUser)
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
		// userRoutes.POST("/picture", userController.SaveFile)
	}

	lobbyRoutes := r.Group("api/lobby", middleware.AuthorizeJWT(jwtService))
	{
		lobbyRoutes.POST("/create", lobbyController.Create)
		lobbyRoutes.GET("/join/:lobbyId", lobbyController.Join)
		lobbyRoutes.GET("/leave/:lobbyId", lobbyController.Leave)
		// userRoutes.POST("/picture", userController.SaveFile)
	}

	gameRoutes := r.Group("api/game", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		gameRoutes.GET("", gameController.GetGame)
		gameRoutes.POST("/start/:lobbyId", gameController.Start)
		gameRoutes.GET("/roll-dice/:dice", gameController.RollDice)
		gameRoutes.GET("/re-roll-dice", gameController.RollDice)
		gameRoutes.GET("/change-turn", gameController.ChangeTurn)
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
