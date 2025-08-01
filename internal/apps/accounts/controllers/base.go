package accounts

import (
	services "lemfi/simplebank/internal/apps/accounts/services"
)

type AccountController struct {
	accountService services.AccountServiceInterface
}

func NewAccountController() *AccountController {
	return &AccountController{
		accountService: services.NewAccountService(),
	}
}
