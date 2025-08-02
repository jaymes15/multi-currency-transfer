package currencies

type Currency string

const (
	USD Currency = "USD"
	NGN Currency = "NGN"
	GBP Currency = "GBP"
	EUR Currency = "EUR"
)

var SupportedCurrencies = []Currency{USD, NGN, GBP, EUR}

// GetSupportedCurrenciesString returns supported currencies as a comma-separated string
func GetSupportedCurrenciesString() string {
	currencies := SupportedCurrencies
	result := ""
	for i, currency := range currencies {
		if i > 0 {
			result += ", "
		}
		result += string(currency)
	}
	return result
}

// GetSupportedCurrenciesForValidation returns supported currencies as a space-separated string for validation tags
func GetSupportedCurrenciesForValidation() string {
	currencies := SupportedCurrencies
	result := ""
	for i, currency := range currencies {
		if i > 0 {
			result += " "
		}
		result += string(currency)
	}
	return result
}
