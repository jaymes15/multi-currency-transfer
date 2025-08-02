package transfers

import "errors"

// Custom error types with status codes
type ClientError struct {
	Message string
	Status  int
}

func (e ClientError) Error() string {
	return e.Message
}

// Predefined client errors
var (
	ErrSameAccountTransfer = ClientError{
		Message: "cannot transfer to the same account",
		Status:  400,
	}
	ErrInvalidAmount = ClientError{
		Message: "transfer amount must be positive",
		Status:  400,
	}
	ErrFromAccountNotFound = ClientError{
		Message: "from account not found",
		Status:  400,
	}
	ErrToAccountNotFound = ClientError{
		Message: "to account not found",
		Status:  400,
	}
	ErrFromAccountCurrencyMismatch = ClientError{
		Message: "from account currency mismatch",
		Status:  400,
	}
	ErrToAccountCurrencyMismatch = ClientError{
		Message: "to account currency mismatch",
		Status:  400,
	}
	ErrInsufficientBalance = ClientError{
		Message: "insufficient balance",
		Status:  400,
	}
)

// Helper function to check if error is a ClientError
func IsClientError(err error) (ClientError, bool) {
	var clientErr ClientError
	if errors.As(err, &clientErr) {
		return clientErr, true
	}
	return ClientError{}, false
}
