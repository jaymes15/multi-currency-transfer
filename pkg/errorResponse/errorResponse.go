package errorResponse

import (
	"fmt"
	"lemfi/simplebank/config"
	helpers "lemfi/simplebank/pkg/helpers"
	"net/http"
	"runtime/debug"
)

// The logError() method is a generic helper for logging an error message. Later in the
// book we'll upgrade this to use structured logging, and record additional information
// about the request including the HTTP method and URL.
func logError(r *http.Request, err error) {
	stackTrace := debug.Stack()
	config.Logger.Error("An error occurred",
		"error", err.Error(),
		"stackTrace", string(stackTrace),
		"method", r.Method,
		"url", r.URL.String(),
	)
}

// The errorResponse() method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code. Note that we're using an interface{}
// type for the message parameter, rather than just a string type, as this gives us
// more flexibility over the values that we can include in the response.
func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := helpers.Envelope{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := helpers.WriteJSON(w, status, env, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(500)
	}
}

// The serverErrorResponse() method will be used when our application encounters an
// unexpected problem at runtime. It logs the detailed error message, then uses the
// errorResponse() helper to send a 500 Internal Server Error status code and JSON
// response (containing a generic error message) to the client.
func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

// The notFoundResponse() method will be used to send a 404 Not Found status code and
// JSON response to the client.
func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}

// The methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func UnAuthorizedRequestResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid or missing authentication token"
	errorResponse(w, r, http.StatusUnauthorized, message)
}
