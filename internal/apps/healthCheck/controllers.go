package healthcheck

import (
	"lemfi/simplebank/config"
	"lemfi/simplebank/pkg/errorResponse"
	"lemfi/simplebank/pkg/responseHandler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheckHandler(c *gin.Context) {
	configs := config.Get()

	// Extract request headers
	headers := map[string]string{}
	for name, values := range c.Request.Header {
		if len(values) > 0 {
			headers[name] = values[0] // Only store the first value for each header
		}
	}

	data := responseHandler.Envelope{
		"status": "available",
		"system_info": map[string]interface{}{
			"environment": configs.Env,
			"headers":     headers, // Include request headers in the response
		},
	}

	err := responseHandler.WriteJSON(c.Writer, http.StatusOK, data, nil)
	if err != nil {
		errorResponse.ServerErrorResponse(c, err)
		return
	}
}
