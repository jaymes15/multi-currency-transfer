package routes

import (
	"lemfi/simplebank/internal/apps/accounts"
	exchangeRates "lemfi/simplebank/internal/apps/exchangeRates"
	healthcheck "lemfi/simplebank/internal/apps/healthCheck"
	transfers "lemfi/simplebank/internal/apps/transfers"
	users "lemfi/simplebank/internal/apps/users"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) *gin.Engine {

	healthcheck.Routes(router)
	accounts.Routes(router)
	transfers.Routes(router)
	exchangeRates.Routes(router)
	users.Routes(router)

	return router
}
