package accounts

import (
	accounts "lemfi/simplebank/internal/apps/accounts/controllers"
	respositories "lemfi/simplebank/internal/apps/accounts/respositories"
	services "lemfi/simplebank/internal/apps/accounts/services"
	"lemfi/simplebank/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Routes defines the health check route for the API.
func Routes(router *gin.Engine) {
	accountRespository := respositories.NewAccountRespository()
	accountService := services.NewAccountService(accountRespository)
	accountController := accounts.NewAccountController(accountService)

	// Group accounts routes with common middleware
	accountsGroup := router.Group("/api/v1/accounts")
	accountsGroup.Use(
		middleware.ValidateAuth(),
		middleware.RequireAuthenticatedUser(),
	)

	// Register routes without repeating middleware
	accountsGroup.POST("", accountController.CreateAccountController)
	accountsGroup.GET("", accountController.GetAccountsController)
}
