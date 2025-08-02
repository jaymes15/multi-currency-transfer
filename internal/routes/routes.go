package routes

import (
	"lemfi/simplebank/internal/apps/accounts"
	healthcheck "lemfi/simplebank/internal/apps/healthCheck"
	transfers "lemfi/simplebank/internal/apps/transfers"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) *gin.Engine {

	healthcheck.Routes(router)
	accounts.Routes(router)
	transfers.Routes(router)

	return router
}
