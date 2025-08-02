package accounts

import (
	"lemfi/simplebank/config"
	requests "lemfi/simplebank/internal/apps/accounts/requests"
	responses "lemfi/simplebank/internal/apps/accounts/responses"
	"lemfi/simplebank/internal/apps/currencies"
)

func (accountService *AccountService) CreateAccount(payload requests.CreateAccountRequest) (responses.CreateAccountResponse, error) {
	config.Logger.Info("Processing account creation in service layer", "owner", payload.Owner, "currency", payload.Currency)

	if !currencies.IsSupportedCurrency(currencies.Currency(payload.Currency)) {
		config.Logger.Error("Currency is not supported", "currency", payload.Currency)
		return responses.CreateAccountResponse{}, currencies.ErrCurrencyNotSupported
	}

	account, err := accountService.accountRespository.CreateAccount(payload)
	if err != nil {
		config.Logger.Error("Failed to create account in service layer", "error", err.Error(), "owner", payload.Owner)
		return responses.CreateAccountResponse{}, err
	}

	config.Logger.Info("Account created successfully in service layer", "accountID", account.ID, "owner", account.Owner)

	response := responses.CreateAccountResponse{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}

	config.Logger.Info("Account creation service completed", "accountID", response.ID)

	return response, nil
}
