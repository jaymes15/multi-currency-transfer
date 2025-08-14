package users

import (
	controllers "lemfi/simplebank/internal/apps/users/controllers"
	respositories "lemfi/simplebank/internal/apps/users/respositories"
	services "lemfi/simplebank/internal/apps/users/services"
	"lemfi/simplebank/pkg/token"

	"github.com/gin-gonic/gin"
)

// Routes defines the user routes for the API.
func Routes(router *gin.Engine) {
	userRespository := respositories.NewUserRespository()
	tokenMaker := token.GetTokenMaker()
	userService := services.NewUserService(userRespository, tokenMaker)
	userController := controllers.NewUserController(userService, tokenMaker)

	router.POST("/api/v1/users", userController.CreateUserController)
	router.POST("/api/v1/users/login", userController.LoginUserController)
}
