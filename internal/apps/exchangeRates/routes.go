package exchangeRates

import (
	controllers "lemfi/simplebank/internal/apps/exchangeRates/controllers"
	respositories "lemfi/simplebank/internal/apps/exchangeRates/respositories"
	services "lemfi/simplebank/internal/apps/exchangeRates/services"

	"github.com/gin-gonic/gin"
)

// Routes defines the exchange rate routes for the API.
func Routes(router *gin.Engine) {
	exchangeRateRepository := respositories.NewExchangeRateRepository()
	exchangeRateService := services.NewExchangeRateService(exchangeRateRepository)
	exchangeRateController := controllers.NewExchangeRateController(exchangeRateService)

	router.GET("/api/v1/exchange-rates", exchangeRateController.ListExchangeRatesController)
	router.POST("/api/v1/exchange-rates/calculate", exchangeRateController.GetExchangeRateController)
}
