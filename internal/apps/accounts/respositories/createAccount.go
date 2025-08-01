package accounts

import (
	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/accounts/requests"
)

func (accountRespository *AccountRespository) CreateAccount(payload requests.CreateAccountRequest) (db.Account, error) {
	config.Logger.Info("Creating account in database", "owner", payload.Owner, "currency", payload.Currency)

	account, err := accountRespository.queries.CreateAccount(accountRespository.context, db.CreateAccountParams{
		Owner:    payload.Owner,
		Balance:  0,
		Currency: payload.Currency,
	})

	if err != nil {
		config.Logger.Error("Failed to create account in database", "error", err.Error(), "owner", payload.Owner)
		return db.Account{}, err
	}

	config.Logger.Info("Successfully created account in database", "accountID", account.ID, "owner", account.Owner, "currency", account.Currency)

	return account, nil
}
