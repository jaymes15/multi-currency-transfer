package exchangeRates

import (
	"context"

	"lemfi/simplebank/config"
	"lemfi/simplebank/internal/apps/currencies"
	exchangeRateErrors "lemfi/simplebank/internal/apps/exchangeRates/errors"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"

	"github.com/shopspring/decimal"
)

func (exchangeRateService *ExchangeRateService) GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (responses.GetExchangeRateResponse, error) {
	config.Logger.Info("Service: Getting exchange rate for currency pair",
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

	dbExchangeRate, err := exchangeRateService.exchangeRateRepository.GetExchangeRate(ctx, payload)
	if err != nil {
		config.Logger.Error("Service: Failed to get exchange rate", "error", err.Error())
		return responses.GetExchangeRateResponse{}, err
	}

	exchangeRate := responses.NewExchangeRateResponse(dbExchangeRate)

	amountToSend := payload.Amount
	amountToReceive := payload.Amount.Mul(dbExchangeRate.Rate).Round(2)

	// Calculate fee based on configuration
	cfg := config.Get()
	fee := cfg.MultiCurrency.Fee
	totalAmount := amountToSend.Add(fee)

	canTransact := false
	message := "Exchange rate expired"
	if !exchangeRateService.IsExchangeRateExpired(dbExchangeRate) {
		canTransact = true
		message = "Exchange rate available for transaction"
	}

	response := responses.GetExchangeRateResponse{
		ExchangeRate:    exchangeRate,
		AmountToSend:    amountToSend,
		AmountToReceive: amountToReceive,
		Fee:             fee,
		TotalAmount:     totalAmount,
		CanTransact:     canTransact,
		Message:         message,
	}

	config.Logger.Info("Successfully fetched exchange rate",
		"rate", response.ExchangeRate.Rate,
		"amount_to_send", amountToSend.String(),
		"amount_to_receive", amountToReceive.String(),
		"fee", fee.String(),
		"total_amount", totalAmount.String(),
	)

	config.Logger.Info("Service: Successfully got exchange rate",
		"rate", response.ExchangeRate.Rate.String(),
		"can_transact", response.CanTransact,
	)

	return response, nil
}
