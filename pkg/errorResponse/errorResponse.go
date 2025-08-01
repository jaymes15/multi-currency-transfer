package errorResponse

import (
	"fmt"
	"lemfi/simplebank/config"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// The logError() method is a generic helper for logging an error message. Later in the
// book we'll upgrade this to use structured logging, and record additional information
// about the request including the HTTP method and URL.
func logError(c *gin.Context, err error) {
	stackTrace := debug.Stack()
	config.Logger.Error("An error occurred",
		"error", err.Error(),
		"stackTrace", string(stackTrace),
		"method", c.Request.Method,
		"url", c.Request.URL.String(),
	)
}

// The errorResponse() method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code. Note that we're using an interface{}
// type for the message parameter, rather than just a string type, as this gives us
// more flexibility over the values that we can include in the response.
func errorResponse(c *gin.Context, status int, message interface{}) {
	c.JSON(status, gin.H{"error": message})
}

// The serverErrorResponse() method will be used when our application encounters an
// unexpected problem at runtime. It logs the detailed error message, then uses the
// errorResponse() helper to send a 500 Internal Server Error status code and JSON
// response (containing a generic error message) to the client.
func ServerErrorResponse(c *gin.Context, err error) {
	logError(c, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(c, http.StatusInternalServerError, message)
}

// The notFoundResponse() method will be used to send a 404 Not Found status code and
// JSON response to the client.
func NotFoundResponse(c *gin.Context) {
	message := "the requested resource could not be found"
	errorResponse(c, http.StatusNotFound, message)
}

// The methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func MethodNotAllowedResponse(c *gin.Context) {
	message := fmt.Sprintf("the %s method is not supported for this resource", c.Request.Method)
	errorResponse(c, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(c *gin.Context, err error) {
	errorResponse(c, http.StatusBadRequest, err.Error())
}

func UnAuthorizedRequestResponse(c *gin.Context) {
	message := "invalid or missing authentication token"
	errorResponse(c, http.StatusUnauthorized, message)
}
