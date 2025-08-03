package exchangeRates

import (
	"context"

	"lemfi/simplebank/config"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"
)

func (exchangeRateService *ExchangeRateService) GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (responses.GetExchangeRateResponse, error) {
	config.Logger.Info("Service: Getting exchange rate for currency pair",
		"from_currency", payload.FromCurrency,
		"to_currency", payload.ToCurrency,
		"amount", payload.Amount.String(),
	)

	dbExchangeRate, err := exchangeRateService.exchangeRateRepository.GetExchangeRate(ctx, payload)
	if err != nil {
		config.Logger.Error("Service: Failed to get exchange rate", "error", err.Error())
		return responses.GetExchangeRateResponse{}, err
	}

	exchangeRate := responses.NewExchangeRateResponse(dbExchangeRate)

	amountToSend := payload.Amount
	amountToReceive := payload.Amount.Mul(dbExchangeRate.Rate).Round(2)

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
		CanTransact:     canTransact,
		Message:         message,
	}

	config.Logger.Info("Successfully fetched exchange rate",
		"rate", response.ExchangeRate.Rate,
		"amount_to_send", amountToSend.String(),
		"amount_to_receive", amountToReceive.String(),
	)

	config.Logger.Info("Service: Successfully got exchange rate",
		"rate", response.ExchangeRate.Rate.String(),
		"can_transact", response.CanTransact,
	)

	return response, nil
}
