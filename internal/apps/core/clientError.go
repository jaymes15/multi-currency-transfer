package core

import "errors"

// Custom error types with status codes
type ClientError struct {
	Message string
	Status  int
}

func (e ClientError) Error() string {
	return e.Message
}

// Helper function to check if error is a ClientError
func IsClientError(err error) (ClientError, bool) {
	var clientErr ClientError
	if errors.As(err, &clientErr) {
		return clientErr, true
	}
	return ClientError{}, false
}
