package transfers

import (
	"lemfi/simplebank/config"
	transferErrors "lemfi/simplebank/internal/apps/transfers/errors"
	requests "lemfi/simplebank/internal/apps/transfers/requests"
	responses "lemfi/simplebank/internal/apps/transfers/responses"
)

func (transferService *TransferService) MakeTransfer(payload requests.MakeTransferRequest) (responses.MakeTransferResponse, error) {
	config.Logger.Info("Processing transfer request",
		"from_account_id", payload.FromAccountID,
		"to_account_id", payload.ToAccountID,
		"amount", payload.Amount,
		"from_currency", payload.FromCurrency,
		"to_currency", payload.ToCurrency,
	)

	// Business validation: Same account transfer prevention
	if payload.FromAccountID == payload.ToAccountID {
		config.Logger.Error("Cannot transfer to same account", "account_id", payload.FromAccountID)
		return responses.MakeTransferResponse{}, transferErrors.ErrSameAccountTransfer
	}

	// Business validation: Amount must be positive
	if payload.Amount <= 0 {
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

	// Execute transfer through repository (includes data validation: account existence, balance check, currency matching)
	result, err := transferService.transferRespository.MakeTransfer(payload)
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
