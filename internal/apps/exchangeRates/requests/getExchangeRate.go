package exchangeRates

import "github.com/shopspring/decimal"

type GetExchangeRateRequest struct {
	FromCurrency string          `json:"from_currency" validate:"required"`
	ToCurrency   string          `json:"to_currency" validate:"required"`
	Amount       decimal.Decimal `json:"amount" validate:"required"`
}
