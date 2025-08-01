package accounts

import (
	accounts "lemfi/simplebank/internal/apps/accounts/controllers"
	respositories "lemfi/simplebank/internal/apps/accounts/respositories"
	services "lemfi/simplebank/internal/apps/accounts/services"

	"github.com/gin-gonic/gin"
)

// Routes defines the health check route for the API.
func Routes(router *gin.Engine) {
	accountRespository := respositories.NewAccountRespository()
	accountService := services.NewAccountService(accountRespository)
	accountController := accounts.NewAccountController(accountService)

	router.POST("/api/v1/accounts", accountController.CreateAccountController)
	router.GET("/api/v1/accounts", accountController.GetAccountsController)

}
