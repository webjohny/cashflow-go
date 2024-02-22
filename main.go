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
	raceRepository   repository.RaceRepository        = repository.NewRaceRepository(db)
	userRepository   repository.UserRepository        = repository.NewUserRepository(db)
	playerRepository repository.PlayerRepository      = repository.NewPlayerRepository(db)
	trxRepository    repository.TransactionRepository = repository.NewTransactionRepository(db)

	// Services
	jwtService         service.JWTService         = service.NewJWTService()
	userService        service.UserService        = service.NewUserService(userRepository)
	transactionService service.TransactionService = service.NewTransactionService(trxRepository)
	playerService      service.PlayerService      = service.NewPlayerService(playerRepository, transactionService)
	authService        service.AuthService        = service.NewAuthService(userRepository)
	gameService        service.GameService        = service.NewGameService(raceRepository)
	raceService        service.RaceService        = service.NewRaceService(raceRepository, playerService, transactionService)
	cardService        service.CardService        = service.NewCardService(gameService, raceService)

	// Controllers
	cardController controller.CardController = controller.NewCardController(cardService)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.Use(middleware.PreRequest())

	cardRoutes := r.Group("api/card")
	{
		//cardRoutes.POST("/:action/:family/:type", cardController.Action)
		cardRoutes.GET("/:action/:family/:type", cardController.Action)
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
