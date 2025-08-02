package currencies

func IsSupportedCurrency(currency Currency) bool {
	switch currency {
	case USD, NGN, GBP, EUR:
		return true
	}
	return false
}
