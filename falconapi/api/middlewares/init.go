package middlewares

import (
	"falconapi/infrastructure/identity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"log"
)

func InitGinMiddlewares(app *gin.Engine,
	initPublicRoutes func(app *gin.Engine),
	initProtectedRoutes func(app *gin.Engine)) {

	logger := logrus.New()
	logger.Level = logrus.TraceLevel
	logrus.SetOutput(gin.DefaultWriter)

	app.Use(gin.Recovery())
	app.Use(CORSMiddleware())
	app.Use(Logger(logger))

	// routes that don't require a JWT token
	initPublicRoutes(app)

	tokenRetrospect := identity.NewIdentityManager()
	app.Use(NewGinJwtMiddleware(tokenRetrospect))

	// routes that require authentication/authorization
	initProtectedRoutes(app)

	log.Println("gin middlewares initialized")
}
