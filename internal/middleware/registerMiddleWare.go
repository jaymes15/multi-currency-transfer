package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterMiddleware(router *gin.Engine) http.Handler {
	// Gin has its own middleware system, no need for alice wrapper
	return router
}
