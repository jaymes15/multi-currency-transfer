package accounts

import (
	accounts "lemfi/simplebank/internal/apps/accounts/controllers"

	"github.com/gin-gonic/gin"
)

// Routes defines the health check route for the API.
func Routes(router *gin.Engine) {
	accountController := accounts.NewAccountController()

	router.POST("/api/v1/accounts", accountController.CreateAccountController)

}
