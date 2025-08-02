package transfers

import (
	"context"
	"lemfi/simplebank/config"
	"lemfi/simplebank/internal/apps/currencies"
	exchangeRateRequests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	transferErrors "lemfi/simplebank/internal/apps/transfers/errors"
	requests "lemfi/simplebank/internal/apps/transfers/requests"
	responses "lemfi/simplebank/internal/apps/transfers/responses"

	"github.com/shopspring/decimal"
)

func (transferService *TransferService) MakeTransfer(payload requests.MakeTransferRequest) (responses.MakeTransferResponse, error) {
	config.Logger.Info("Processing transfer request",
		"from_account_id", payload.FromAccountID,
		"to_account_id", payload.ToAccountID,
		"amount", payload.Amount,
		"from_currency", payload.FromCurrency,
		"to_currency", payload.ToCurrency,
	)

	if !currencies.IsSupportedCurrency(currencies.Currency(payload.FromCurrency)) {
		config.Logger.Error("From currency is not supported", "currency", payload.FromCurrency)
		return responses.MakeTransferResponse{}, currencies.ErrCurrencyNotSupported
	}

	// Business validation: Same account transfer prevention
	if payload.FromAccountID == payload.ToAccountID {
		config.Logger.Error("Cannot transfer to same account", "account_id", payload.FromAccountID)
		return responses.MakeTransferResponse{}, transferErrors.ErrSameAccountTransfer
	}

	// Business validation: Amount must be positive
	if payload.Amount.LessThanOrEqual(decimal.Zero) {
		config.Logger.Error("Invalid transfer amount", "amount", payload.Amount)
		return responses.MakeTransferResponse{}, transferErrors.ErrInvalidAmount
	}

	// Business validation: Currency consistency
	if payload.FromCurrency == payload.ToCurrency {
		config.Logger.Info("Same currency transfer", "currency", payload.FromCurrency)
	} else {
		config.Logger.Info("Cross-currency transfer",
			"from_currency", payload.FromCurrency,
			"to_currency", payload.ToCurrency,
		)
	}

	// Calculate exchange rate and converted amount in service layer
	convertedAmount := payload.Amount
	exchangeRate := decimal.NewFromInt(1) // Default to 1:1 for same currency

	if payload.FromCurrency != payload.ToCurrency {
		// Use exchange rate service to get exchange rate
		exchangeRateRequest := exchangeRateRequests.GetExchangeRateRequest{
			FromCurrency: payload.FromCurrency,
			ToCurrency:   payload.ToCurrency,
			Amount:       payload.Amount,
		}

		exchangeRateResponse, err := transferService.exchangeRateService.GetExchangeRate(context.Background(), exchangeRateRequest)
		if err != nil {
			config.Logger.Error("Failed to get exchange rate",
				"from_currency", payload.FromCurrency,
				"to_currency", payload.ToCurrency,
				"error", err.Error(),
			)
			return responses.MakeTransferResponse{}, err
		}

		exchangeRate = exchangeRateResponse.ExchangeRate.Rate
		convertedAmount = exchangeRateResponse.AmountToReceive
	}

	// Execute transfer through repository (includes data validation: account existence, balance check, currency matching)
	result, err := transferService.transferRespository.MakeTransfer(payload, convertedAmount, exchangeRate)
	if err != nil {
		config.Logger.Error("Transfer failed",
			"error", err.Error(),
			"from_account_id", payload.FromAccountID,
			"to_account_id", payload.ToAccountID,
		)
		return responses.MakeTransferResponse{}, err
	}

	// Convert database result to response using helper function
	response := responses.NewMakeTransferResponse(result)

	config.Logger.Info("Transfer completed successfully",
		"transfer_id", result.Transfer.ID,
		"from_account_id", result.Transfer.FromAccountID,
		"to_account_id", result.Transfer.ToAccountID,
		"amount", result.Transfer.Amount,
		"from_balance", result.FromAccount.Balance,
		"to_balance", result.ToAccount.Balance,
	)

	return response, nil
}
