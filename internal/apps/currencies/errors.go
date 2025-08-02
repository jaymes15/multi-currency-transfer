package currencies

import "lemfi/simplebank/internal/apps/core"

var (
	ErrCurrencyNotSupported = core.ClientError{
		Message: "currency is not supported. Supported currencies are: " + GetSupportedCurrenciesString(),
		Status:  400,
	}
)
