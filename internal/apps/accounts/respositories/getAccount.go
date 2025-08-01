package accounts

import (
	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
)

func (accountRespository *AccountRespository) GetAccounts() ([]db.Account, error) {
	config.Logger.Info("Fetching accounts from database", "limit", 10, "offset", 0)

	accounts, err := accountRespository.queries.ListAllAccounts(accountRespository.context, db.ListAllAccountsParams{
		Limit:  10,
		Offset: 0,
	})

	if err != nil {
		config.Logger.Error("Failed to fetch accounts from database", "error", err.Error())
		return []db.Account{}, err
	}

	config.Logger.Info("Successfully fetched accounts from database", "count", len(accounts))

	return accounts, nil
}
