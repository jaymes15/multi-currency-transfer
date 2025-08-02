package accounts

import (
	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
	accountErrors "lemfi/simplebank/internal/apps/accounts/errors"
	requests "lemfi/simplebank/internal/apps/accounts/requests"
	"strings"

	"github.com/shopspring/decimal"
)

func (accountRespository *AccountRespository) CreateAccount(payload requests.CreateAccountRequest) (db.Account, error) {
	config.Logger.Info("Creating account in database", "owner", payload.Owner, "currency", payload.Currency)

	account, err := accountRespository.queries.CreateAccount(accountRespository.context, db.CreateAccountParams{
		Owner:    payload.Owner,
		Balance:  decimal.Zero,
		Currency: payload.Currency,
	})

	if err != nil {
		// Check if it's a unique constraint violation
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") &&
			strings.Contains(err.Error(), "unique_owner_currency") {
			config.Logger.Error("Duplicate account creation attempted", "owner", payload.Owner, "currency", payload.Currency)
			return db.Account{}, accountErrors.ErrDuplicateAccount
		}

		config.Logger.Error("Failed to create account in database", "error", err.Error(), "owner", payload.Owner)
		return db.Account{}, err
	}

	config.Logger.Info("Successfully created account in database", "accountID", account.ID, "owner", account.Owner, "currency", account.Currency)

	return account, nil
}
