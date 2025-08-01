package accounts

import (
	services "lemfi/simplebank/internal/apps/accounts/services"
)

type AccountController struct {
	accountService services.AccountServiceInterface
}

func NewAccountController(service services.AccountServiceInterface) *AccountController {
	return &AccountController{
		accountService: service,
	}
}
