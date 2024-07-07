package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
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
	usedCardRepository   repository.UsedCardRepository    = repository.NewUsedCardRepository(db)
	playerRepository     repository.PlayerRepository      = repository.NewPlayerRepository(db)
	professionRepository repository.ProfessionRepository  = repository.NewProfessionRepository(os.Getenv("PROFESSIONS_PATH"))
	trxRepository        repository.TransactionRepository = repository.NewTransactionRepository(db)

	// Services
	jwtService         service.JWTService         = service.NewJWTService()
	userService        service.UserService        = service.NewUserService(userRepository)
	transactionService service.TransactionService = service.NewTransactionService(trxRepository)
	professionService  service.ProfessionService  = service.NewProfessionService(professionRepository)
	playerService      service.PlayerService      = service.NewPlayerService(playerRepository, professionService, transactionService)
	authService        service.AuthService        = service.NewAuthService(userRepository)
	gameService        service.GameService        = service.NewGameService(raceService, playerService, lobbyService, professionService)
	raceService        service.RaceService        = service.NewRaceService(raceRepository, playerService, transactionService)
	lobbyService       service.LobbyService       = service.NewLobbyService(lobbyRepository)
	cardService        service.CardService        = service.NewCardService(usedCardRepository, gameService, raceService, playerService)
	financeService     service.FinanceService     = service.NewFinanceService(raceService, playerService)

	// Controllers
	gameController       controller.GameController       = controller.NewGameController(gameService)
	playerController     controller.PlayerController     = controller.NewPlayerController(playerService)
	playerTestController controller.PlayerTestController = controller.NewPlayerTestController(playerService)
	lobbyController      controller.LobbyController      = controller.NewLobbyController(lobbyService)
	financeController    controller.FinanceController    = controller.NewFinanceController(financeService)
	cardController       controller.CardController       = controller.NewCardController(cardService)
	authController       controller.AuthController       = controller.NewAuthController(authService, jwtService)
	userController       controller.UserController       = controller.NewUserController(userService, jwtService)
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
			"error": storage.ErrorForbidden,
		})
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	//log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
	log.SetOutput(colorable.NewColorableStdout())
	log.SetReportCaller(true)

	log.SetOutput(os.Stdout)

	//log.SetLevel(log.WarnLevel)
}

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()
	//gin.ReleaseMode
	gin.SetMode("debug")

	cardRoutes := r.Group("api/card", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		cardRoutes.GET("/test/:action", cardController.TestCard)

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
		lobbyRoutes.GET("/:lobbyId", lobbyController.GetLobby)
		lobbyRoutes.POST("/create", lobbyController.Create)
		lobbyRoutes.GET("/join/:lobbyId", lobbyController.Join)
		lobbyRoutes.GET("/leave/:lobbyId", lobbyController.Leave)
		lobbyRoutes.GET("/cancel/:lobbyId", lobbyController.Cancel)
		// userRoutes.POST("/picture", userController.SaveFile)
	}

	gameRoutes := r.Group("api/game", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		gameRoutes.GET("/:raceId", gameController.GetGame)
		gameRoutes.GET("/cancel/:raceId", gameController.Cancel)
		gameRoutes.GET("/reset/:raceId", gameController.Reset)
		gameRoutes.POST("/start/:lobbyId", gameController.Start)
		// @toDo Create EP for moving to big race
		gameRoutes.POST("/move/:raceId", gameController.MoveToBigRace)
		gameRoutes.GET("/roll-dice/:dice", gameController.RollDice)
		gameRoutes.GET("/re-roll-dice", gameController.RollDice)
		gameRoutes.GET("/change-turn", gameController.ChangeTurn)
		gameRoutes.GET("/get/tiles", gameController.GetTiles)
	}

	financeRoutes := r.Group("api/finance", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		financeRoutes.POST("/send/money", financeController.SendMoney)
		financeRoutes.POST("/send/assets", financeController.SendAssets)
		financeRoutes.POST("/loan/take", financeController.TakeLoan)
	}

	playerRoutes := r.Group("api/player", middleware.AuthorizeJWT(jwtService), middleware.GetGameId())
	{
		playerRoutes.GET("/info", playerController.GetRacePlayer)
	}

	playerTestRoutes := r.Group("test/player")
	{
		playerTestRoutes.GET("/", playerTestController.Index)
		playerTestRoutes.GET("/sell-stocks", playerTestController.SellStocks)

		playerTestRoutes.GET("/increase-stocks", playerTestController.IncreaseStocks)
		playerTestRoutes.GET("/decrease-stocks", playerTestController.DecreaseStocks)

		playerTestRoutes.GET("/buy-stocks", playerTestController.BuyStocks)
		playerTestRoutes.GET("/buy-other-assets", playerTestController.BuyOtherAssets)
		playerTestRoutes.GET("/buy-real-estate", playerTestController.BuyRealEstate)
		playerTestRoutes.GET("/buy-lottery", playerTestController.BuyLottery)
		playerTestRoutes.GET("/buy-business", playerTestController.BuyBusiness)
		playerTestRoutes.GET("/buy-partner-real-estate", playerTestController.BuyRealEstateInPartnership)
		playerTestRoutes.GET("/buy-partner-business", playerTestController.BuyBusinessInPartnership)

		playerTestRoutes.GET("/buy-dream", playerTestController.BuyDream)
		playerTestRoutes.GET("/buy-big-business", playerTestController.BuyBigBusiness)
		playerTestRoutes.GET("/buy-risk-business", playerTestController.BuyRiskBusiness)
		playerTestRoutes.GET("/buy-risk-stocks", playerTestController.BuyRiskStocks)
	}

	r.Run(":" + os.Getenv("PORT"))
}
