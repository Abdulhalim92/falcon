package routes

import (
	"falconapi/api/handlers"
	"falconapi/api/middlewares"
	_ "falconapi/docs"
	"falconapi/infrastructure/database"
	"falconapi/infrastructure/identity"
	gorm_db "falconapi/shared/gorm-db"
	"falconapi/use_cases/productsuc"
	"falconapi/use_cases/terminaluc"
	"falconapi/use_cases/usermgmtuc"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

func InitPublicRoutes(app *gin.Engine) {

	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome to My Falcon Api")
	})

	grp := app.Group("/api/v1")

	identityManager := identity.NewIdentityManager()
	registerUseCase := usermgmtuc.NewRegisterUseCase(identityManager)
	loginUseCase := usermgmtuc.NewLoginUseCase(identityManager)
	otpUseCase := usermgmtuc.NewOtpUseCase(identityManager)

	grp.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	grp.POST("/user", handlers.Register(registerUseCase))
	grp.POST("/login", handlers.Login(loginUseCase))
	grp.POST("/generate-otp", handlers.GenerateOtp(otpUseCase))
	grp.POST("/validate-otp", handlers.ValidateOtp(otpUseCase))
}

func InitProtectedRoutes(app *gin.Engine) {

	grp := app.Group("/api/v1")

	db, err := gorm_db.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	gorm_db.MigrateDatabase(db)

	//productsDataStore := datastores.NewProductsDataStore()
	productsDataStore := database.NewProductDatabase(db)
	terminalDataStore := database.NewTerminalDataStore(db)

	createProductUseCase := productsuc.NewCreateProductUseCase(productsDataStore)
	grp.POST("/products", middlewares.NewRequiresRealmRole("admin"),
		handlers.CreateProductHandler(createProductUseCase))

	getProductsUseCase := productsuc.NewGetProductsUseCase(productsDataStore)
	grp.GET("/products", middlewares.NewRequiresRealmRole("viewer"),
		handlers.GetProductsHandler(getProductsUseCase))

	getTerminalStatusesUseCase := terminaluc.NewGetTerminalStatusesUseCase(terminalDataStore)
	grp.GET("/terminalstatuses", handlers.GetTerminalsStatusesHandler(getTerminalStatusesUseCase))

	getTerminalInfoUseCase := terminaluc.NewGetTerminalInfoUseCase(terminalDataStore)
	grp.GET("/terminalsinfo", handlers.GetTerminalsInfoHandler(getTerminalInfoUseCase))

	getRegionsUseCase := terminaluc.NewGetRegionsUseCase(terminalDataStore)
	grp.GET("/region", handlers.GetRegionsHandler(getRegionsUseCase))
}
