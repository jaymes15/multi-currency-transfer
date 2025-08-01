package accounts

import (
	requests "lemfi/simplebank/internal/apps/accounts/requests"
	responses "lemfi/simplebank/internal/apps/accounts/responses"
)

type AccountServiceInterface interface {
	CreateAccount(payload requests.CreateAccountRequest) (responses.CreateAccountResponse, error)
	GetAccounts() ([]responses.GetAccountResponse, error)
}
