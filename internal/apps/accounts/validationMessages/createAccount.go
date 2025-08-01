package accounts

var CreateAccountValidationMessages = map[string]string{
	"Owner.required":    "owner is required.",
	"Currency.required": "currency is required.",
	"Currency.oneof":    "currency must be one of the following: USD, NGN, GBP, EUR.",
}
