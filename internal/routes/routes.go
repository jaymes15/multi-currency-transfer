package routes

import (
	"lemfi/simplebank/internal/apps/accounts"
	healthcheck "lemfi/simplebank/internal/apps/healthCheck"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) *gin.Engine {

	healthcheck.Routes(router)
	accounts.Routes(router)

	return router
}
