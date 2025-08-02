package transfers

import (
	controllers "lemfi/simplebank/internal/apps/transfers/controllers"
	respositories "lemfi/simplebank/internal/apps/transfers/respositories"
	services "lemfi/simplebank/internal/apps/transfers/services"

	"github.com/gin-gonic/gin"
)

// Routes defines the health check route for the API.
func Routes(router *gin.Engine) {
	transferRespository := respositories.NewTransferRespository()
	transferService := services.NewTransferService(transferRespository)
	transferController := controllers.NewTransferController(transferService)

	router.POST("/api/v1/transfers", transferController.MakeTransferController)
}
