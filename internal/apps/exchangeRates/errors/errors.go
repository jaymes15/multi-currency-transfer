package errors

import "lemfi/simplebank/internal/apps/core"

var (
	ErrExchangeRateNotFound = core.ClientError{
		Message: "exchange rate not found for currency pair",
		Status:  404,
	}

	ErrInvalidCurrencyPair = core.ClientError{
		Message: "invalid currency pair",
		Status:  400,
	}

	ErrUnsupportedCurrency = core.ClientError{
		Message: "unsupported currency",
		Status:  400,
	}

	ErrInvalidAmount = core.ClientError{
		Message: "invalid amount",
		Status:  400,
	}

	ErrExchangeRateExpired = core.ClientError{
		Message: "exchange rate expired",
		Status:  400,
	}

	ErrExchangeRateMismatch = core.ClientError{
		Message: "exchange rate mismatch",
		Status:  400,
	}

	ErrExchangeRateZero = core.ClientError{
		Message: "exchange rate is zero",
		Status:  400,
	}
)
