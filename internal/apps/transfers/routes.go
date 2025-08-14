package transfers

import (
	exchangeRateRespositories "lemfi/simplebank/internal/apps/exchangeRates/respositories"
	exchangeRateServices "lemfi/simplebank/internal/apps/exchangeRates/services"
	controllers "lemfi/simplebank/internal/apps/transfers/controllers"
	respositories "lemfi/simplebank/internal/apps/transfers/respositories"
	services "lemfi/simplebank/internal/apps/transfers/services"
	"lemfi/simplebank/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Routes defines the health check route for the API.
func Routes(router *gin.Engine) {
	// Initialize repositories
	transferRespository := respositories.NewTransferRespository()
	exchangeRateRepository := exchangeRateRespositories.NewExchangeRateRepository()

	// Initialize services
	exchangeRateService := exchangeRateServices.NewExchangeRateService(exchangeRateRepository)
	transferService := services.NewTransferService(transferRespository, exchangeRateService)

	// Initialize controllers
	transferController := controllers.NewTransferController(transferService)

	// Group transfers routes with common middleware
	transfersGroup := router.Group("/api/v1/transfers")
	transfersGroup.Use(
		middleware.ValidateAuth(),
		middleware.RequireAuthenticatedUser(),
	)

	// Register routes
	transfersGroup.POST("", transferController.MakeTransferController)
}
