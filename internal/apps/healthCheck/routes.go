package healthcheck

import (
	"github.com/gin-gonic/gin"
)

// Routes defines the health check route for the API.
func Routes(router *gin.Engine) {
	router.GET("/api/v1/healthz", HealthCheckHandler)
}
