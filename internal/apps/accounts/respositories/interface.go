package accounts

import (
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/accounts/requests"
)

type AccountRespositoryInterface interface {
	CreateAccount(payload requests.CreateAccountRequest) (db.Account, error)
	GetAccounts() ([]db.Account, error)
}
