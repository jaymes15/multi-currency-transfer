package accounts

import (
	"lemfi/simplebank/config"
	responses "lemfi/simplebank/internal/apps/accounts/responses"
)

func (accountService *AccountService) GetAccounts() ([]responses.GetAccountResponse, error) {
	config.Logger.Info("Processing get accounts request in service layer")

	accounts, err := accountService.accountRespository.GetAccounts()
	if err != nil {
		config.Logger.Error("Failed to get accounts in service layer", "error", err.Error())
		return []responses.GetAccountResponse{}, err
	}

	config.Logger.Info("Successfully retrieved accounts from repository", "count", len(accounts))

	accountsResponse := make([]responses.GetAccountResponse, len(accounts))
	for i, account := range accounts {
		accountsResponse[i] = responses.GetAccountResponse{
			ID:        account.ID,
			Owner:     account.Owner,
			Balance:   account.Balance,
			Currency:  account.Currency,
			CreatedAt: account.CreatedAt,
		}
	}

	config.Logger.Info("Get accounts service completed successfully", "count", len(accountsResponse))

	return accountsResponse, nil
}
