package users

import (
	controllers "lemfi/simplebank/internal/apps/users/controllers"
	respositories "lemfi/simplebank/internal/apps/users/respositories"
	services "lemfi/simplebank/internal/apps/users/services"
	"lemfi/simplebank/internal/middleware"
	"lemfi/simplebank/pkg/token"

	"github.com/gin-gonic/gin"
)

// Routes defines the user routes for the API.
func Routes(router *gin.Engine) {
	userRespository := respositories.NewUserRespository()
	tokenMaker := token.GetTokenMaker()
	userService := services.NewUserService(userRespository, tokenMaker)
	userController := controllers.NewUserController(userService, tokenMaker)

	// Public routes (no authentication required)
	router.POST("/api/v1/users", userController.CreateUserController)
	router.POST("/api/v1/users/login", userController.LoginUserController)
	router.POST("/api/v1/users/refresh", userController.RefreshTokenController)
	router.POST("/api/v1/users/logout", userController.LogoutController)

	// Protected routes (authentication required)
	router.GET("/api/v1/users/me", middleware.ValidateAuth(), middleware.RequireAuthenticatedUser(), userController.GetUserController)
}
