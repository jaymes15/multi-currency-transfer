package errors

import "errors"

type ClientError struct {
	Message    string
	StatusCode int
}

func (e ClientError) Error() string {
	return e.Message
}

func (e ClientError) GetStatusCode() int {
	return e.StatusCode
}

var (
	ErrExchangeRateNotFound = &ClientError{
		Message:    "exchange rate not found for currency pair",
		StatusCode: 404,
	}

	ErrInvalidCurrencyPair = &ClientError{
		Message:    "invalid currency pair",
		StatusCode: 400,
	}

	ErrUnsupportedCurrency = &ClientError{
		Message:    "unsupported currency",
		StatusCode: 400,
	}

	ErrInvalidAmount = &ClientError{
		Message:    "invalid amount",
		StatusCode: 400,
	}
)

// Helper function to check if an error is a client error
func IsClientError(err error) bool {
	var clientErr *ClientError
	return errors.As(err, &clientErr)
}
