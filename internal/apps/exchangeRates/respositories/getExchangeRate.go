package exchangeRates

import (
	"context"

	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
	"lemfi/simplebank/internal/apps/currencies"
	exchangeRateErrors "lemfi/simplebank/internal/apps/exchangeRates/errors"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"

	"github.com/shopspring/decimal"
)

func (exchangeRateRepository *ExchangeRateRepository) GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (responses.GetExchangeRateResponse, error) {
	config.Logger.Info("Fetching exchange rate for currency pair",
		"from_currency", payload.FromCurrency,
		"to_currency", payload.ToCurrency,
		"amount", payload.Amount.String(),
	)

	// Validate currencies are supported
	if !currencies.IsSupportedCurrency(currencies.Currency(payload.FromCurrency)) {
		config.Logger.Error("From currency is not supported", "currency", payload.FromCurrency)
		return responses.GetExchangeRateResponse{}, exchangeRateErrors.ErrUnsupportedCurrency
	}

	if !currencies.IsSupportedCurrency(currencies.Currency(payload.ToCurrency)) {
		config.Logger.Error("To currency is not supported", "currency", payload.ToCurrency)
		return responses.GetExchangeRateResponse{}, exchangeRateErrors.ErrUnsupportedCurrency
	}

	// Validate amount is positive
	if payload.Amount.LessThanOrEqual(decimal.Zero) {
		config.Logger.Error("Invalid amount", "amount", payload.Amount.String())
		return responses.GetExchangeRateResponse{}, exchangeRateErrors.ErrInvalidAmount
	}

	// Get exchange rate from database
	dbExchangeRate, err := exchangeRateRepository.queries.GetExchangeRate(ctx, db.GetExchangeRateParams{
		FromCurrency: payload.FromCurrency,
		ToCurrency:   payload.ToCurrency,
	})
	if err != nil {
		config.Logger.Error("Failed to get exchange rate",
			"from_currency", payload.FromCurrency,
			"to_currency", payload.ToCurrency,
			"error", err.Error(),
		)
		return responses.GetExchangeRateResponse{}, exchangeRateErrors.ErrExchangeRateNotFound
	}

	// Convert database result to response format
	exchangeRate := responses.ExchangeRateResponse{
		ID:           dbExchangeRate.ID,
		FromCurrency: dbExchangeRate.FromCurrency,
		ToCurrency:   dbExchangeRate.ToCurrency,
		Rate:         dbExchangeRate.Rate,
		CreatedAt:    dbExchangeRate.CreatedAt.Time,
	}

	// Calculate amounts
	amountToSend := payload.Amount
	amountToReceive := payload.Amount.Mul(dbExchangeRate.Rate).Round(2)

	// Determine if transaction is possible
	canTransact := true
	message := "Exchange rate available for transaction"

	response := responses.GetExchangeRateResponse{
		ExchangeRate:    exchangeRate,
		AmountToSend:    amountToSend,
		AmountToReceive: amountToReceive,
		CanTransact:     canTransact,
		Message:         message,
	}

	config.Logger.Info("Successfully fetched exchange rate",
		"rate", dbExchangeRate.Rate.String(),
		"amount_to_send", amountToSend.String(),
		"amount_to_receive", amountToReceive.String(),
	)

	return response, nil
}
