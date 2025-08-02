package transfers

var MakeTransferValidationMessages = map[string]string{
	"FromAccountID.required": "from_account_id is required",
	"FromAccountID.min":      "from_account_id must be greater than 0",
	"ToAccountID.required":   "to_account_id is required",
	"ToAccountID.min":        "to_account_id must be greater than 0",
	"Amount.required":        "amount is required",
	"Amount.min":             "amount must be greater than 0",
	"FromCurrency.required":  "from_currency is required",
	"ToCurrency.required":    "to_currency is required",
}
