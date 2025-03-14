package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/config"
	"github.com/webjohny/cashflow-go/controller"
	"github.com/webjohny/cashflow-go/middleware"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"gorm.io/gorm"
	"os"
)

var (
	db *gorm.DB = config.SetupDatabaseConnection()

	// Repositories
	raceRepository        repository.RaceRepository        = repository.NewRaceRepository(db)
	lobbyRepository       repository.LobbyRepository       = repository.NewLobbyRepository(db)
	userRepository        repository.UserRepository        = repository.NewUserRepository(db)
	userRequestRepository repository.UserRequestRepository = repository.NewUserRequestRepository(db)
	playerRepository      repository.PlayerRepository      = repository.NewPlayerRepository(db)
	professionRepository  repository.ProfessionRepository  = repository.NewProfessionRepository(os.Getenv("PROFESSIONS_PATH"))
	trxRepository         repository.TransactionRepository = repository.NewTransactionRepository(db)

	// Services
	jwtService         service.JWTService         = service.NewJWTService()
	userService        service.UserService        = service.NewUserService(userRepository)
	transactionService service.TransactionService = service.NewTransactionService(trxRepository)
	professionService  service.ProfessionService  = service.NewProfessionService(professionRepository)
	playerService      service.PlayerService      = service.NewPlayerService(playerRepository, professionService, transactionService)
	userRequestService service.UserRequestService = service.NewUserRequestService(userRequestRepository)
	authService        service.AuthService        = service.NewAuthService(userRepository)
	gameService        service.GameService        = service.NewGameService(raceService, playerService, lobbyService, professionService)
	raceService        service.RaceService        = service.NewRaceService(raceRepository, playerService, transactionService)
	lobbyService       service.LobbyService       = service.NewLobbyService(lobbyRepository)
	cardService        service.CardService        = service.NewCardService(gameService, raceService, playerService)
	financeService     service.FinanceService     = service.NewFinanceService(userRequestRepository, cardService, raceService, playerService)

	// Controllers
	backdoorController   controller.BackdoorController   = controller.NewBackdoorController(cardService, raceService, playerService, gameService)
	gameController       controller.GameController       = controller.NewGameController(gameService, professionService)
	moderatorController  controller.ModeratorController  = controller.NewModeratorController(playerService, raceService, lobbyService, userRequestService)
	playerController     controller.PlayerController     = controller.NewPlayerController(playerService, raceService, lobbyService)
	playerTestController controller.PlayerTestController = controller.NewPlayerTestController(playerService)
	lobbyController      controller.LobbyController      = controller.NewLobbyController(lobbyService)
	financeController    controller.FinanceController    = controller.NewFinanceController(financeService)
	cardController       controller.CardController       = controller.NewCardController(cardService)
	authController       controller.AuthController       = controller.NewAuthController(authService, jwtService)
	userController       controller.UserController       = controller.NewUserController(userService, jwtService)
)

func init() {
	//log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: false})
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
	log.SetOutput(colorable.NewColorableStdout())
	//log.SetReportCaller(true)

	log.SetOutput(os.Stdout)

	//log.SetLevel(log.WarnLevel)
}

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	//gin.ReleaseMode
	gin.SetMode("debug")

	r.GET("/health", func(ctx *gin.Context) {
		request.FinalResponse(ctx, nil, nil)
	})

	cardRoutes := r.Group("api/card", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		cardRoutes.GET("/test/:action", cardController.TestCard)

		cardRoutes.GET("/type", cardController.Type)
		cardRoutes.GET("/reset-transaction", cardController.ResetTransaction)
		cardRoutes.POST("/cards", cardController.SetCards)
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

	backdoorRoutes := r.Group("api/backdoor", middleware.GetGameId())
	{
		backdoorRoutes.POST("/:raceId/change-card", backdoorController.ChangeCard)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
		// userRoutes.POST("/picture", userController.SaveFile)
	}

	moderatorRoutes := r.Group("api/moderator", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		moderatorRoutes.GET("/:raceId/race", moderatorController.GetRace)
		moderatorRoutes.PUT("/:raceId/status", moderatorController.UpdateStatusRace)
		moderatorRoutes.GET("/:raceId/player", moderatorController.GetRacePlayer)
		moderatorRoutes.POST("/:raceId/send-money", moderatorController.SendMoney)
		moderatorRoutes.GET("/:raceId/players", moderatorController.GetRacePlayers)
		moderatorRoutes.PUT("/:raceId/player/:playerId", moderatorController.UpdatePlayer)
		moderatorRoutes.PUT("/:raceId/race", moderatorController.UpdateRace)
		moderatorRoutes.PUT("/:raceId/handle/user-request", moderatorController.HandleUserRequest)
		// userRoutes.POST("/picture", userController.SaveFile)
	}

	lobbyRoutes := r.Group("api/lobby", middleware.AuthorizeJWT(jwtService))
	{
		lobbyRoutes.GET("/:lobbyId", lobbyController.GetLobby)
		lobbyRoutes.POST("/create", lobbyController.Create)
		lobbyRoutes.GET("/join/:lobbyId", lobbyController.Join)
		lobbyRoutes.GET("/leave/:lobbyId", lobbyController.Leave)
		lobbyRoutes.GET("/cancel/:lobbyId", lobbyController.Cancel)
		lobbyRoutes.PUT("/options/:lobbyId", lobbyController.SetOptions)
		// userRoutes.POST("/picture", userController.SaveFile)
	}

	gameRoutes := r.Group("api/game", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		gameRoutes.GET("/:raceId", gameController.GetGame)
		gameRoutes.GET("/cancel/:raceId", gameController.Cancel)
		gameRoutes.GET("/reset/:raceId", gameController.Reset)
		gameRoutes.POST("/start/:lobbyId", gameController.Start)
		gameRoutes.POST("/roll-dice", gameController.RollDice)
		gameRoutes.GET("/change-turn", gameController.ChangeTurn)
		gameRoutes.GET("/get/tiles", gameController.GetTiles)
	}

	financeRoutes := r.Group("api/finance", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		financeRoutes.POST("/send/money", financeController.SendMoney)
		financeRoutes.POST("/send/assets", financeController.SendAssets)
		financeRoutes.POST("/loan/take", financeController.TakeLoan)
		financeRoutes.POST("/ask/money", financeController.AskMoney)
	}

	playerRoutes := r.Group("api/player", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		playerRoutes.GET("/info", playerController.GetRacePlayer)
		playerRoutes.POST("/on-big-race/:raceId", playerController.MoveOnBigRace)
		playerRoutes.POST("/dream/:raceId", playerController.SetDream)
		playerRoutes.GET("/data/:raceId", playerController.GetPlayerData)
		playerRoutes.PUT("/data/:raceId", playerController.SetPlayerData)
		playerRoutes.POST("/moderator/:raceId", playerController.BecomeModerator)
		playerRoutes.POST("/read-notification/:notificationId/:raceId", playerController.IsReadNotification)
	}

	playerTestRoutes := r.Group("test/player")
	{
		playerTestRoutes.GET("/", playerTestController.Index)
		playerTestRoutes.GET("/extra-money", playerTestController.AddMoney)
		playerTestRoutes.GET("/sell-stocks", playerTestController.SellStocks)
		playerTestRoutes.GET("/sell-business", playerTestController.SellBusiness)
		playerTestRoutes.GET("/sell-real-estate", playerTestController.SellRealEstate)
		playerTestRoutes.GET("/sell-other-assets", playerTestController.SellOtherAssets)

		playerTestRoutes.GET("/damage-real-estate", playerTestController.DamageRealEstate)
		playerTestRoutes.GET("/increase-stocks", playerTestController.IncreaseStocks)
		playerTestRoutes.GET("/decrease-stocks", playerTestController.DecreaseStocks)

		playerTestRoutes.GET("/buy-stocks", playerTestController.BuyStocks)
		playerTestRoutes.GET("/buy-other-assets", playerTestController.BuyOtherAssets)
		playerTestRoutes.GET("/buy-real-estate", playerTestController.BuyRealEstate)
		playerTestRoutes.GET("/buy-lottery", playerTestController.BuyLottery)
		playerTestRoutes.GET("/buy-business", playerTestController.BuyBusiness)
		playerTestRoutes.GET("/buy-partner-other-assets", playerTestController.BuyOtherAssetsInPartnership)
		playerTestRoutes.GET("/buy-partner-real-estate", playerTestController.BuyRealEstateInPartnership)
		playerTestRoutes.GET("/buy-partner-business", playerTestController.BuyBusinessInPartnership)

		playerTestRoutes.GET("/buy-dream", playerTestController.BuyDream)
		playerTestRoutes.GET("/buy-big-business", playerTestController.BuyBigBusiness)
		playerTestRoutes.GET("/buy-risk-business", playerTestController.BuyRiskBusiness)
		playerTestRoutes.GET("/buy-risk-stocks", playerTestController.BuyRiskStocks)
	}

	r.Run(":" + os.Getenv("PORT"))
}
