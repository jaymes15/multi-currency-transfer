package users

import (
	controllers "lemfi/simplebank/internal/apps/users/controllers"
	respositories "lemfi/simplebank/internal/apps/users/respositories"
	services "lemfi/simplebank/internal/apps/users/services"

	"github.com/gin-gonic/gin"
)

// Routes defines the user routes for the API.
func Routes(router *gin.Engine) {
	userRespository := respositories.NewUserRespository()
	userService := services.NewUserService(userRespository)
	userController := controllers.NewUserController(userService)

	router.POST("/api/v1/users", userController.CreateUserController)
} 