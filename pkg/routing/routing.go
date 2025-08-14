package routing

import (
	"net/http"

	"lemfi/simplebank/config"
	"lemfi/simplebank/internal/middleware"
	"lemfi/simplebank/internal/routes"
	"lemfi/simplebank/pkg/errorResponse"

	"github.com/gin-gonic/gin"
)

func Init() {
	config.Logger.Info("Starting Init httprouter....")
	router = gin.Default()
	config.Logger.Info("Completed Init httprouter....")
}

func getRouter() *gin.Engine {
	config.Logger.Info("Returning router *gin.Engine")
	return router
}

func registerRoutes(router *gin.Engine) http.Handler {
	config.Logger.Info("Starting to register all routes")
	router.NoRoute(errorResponse.NotFoundResponse)
	router.NoMethod(errorResponse.MethodNotAllowedResponse)
	registedRoutes := middleware.RegisterMiddleware(routes.Routes(router))

	config.Logger.Info("Completed to register all routes")

	return registedRoutes

}
