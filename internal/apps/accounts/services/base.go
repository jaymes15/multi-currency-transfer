package accounts

import (
	respositories "lemfi/simplebank/internal/apps/accounts/respositories"
)

type AccountService struct {
	accountRespository respositories.AccountRespositoryInterface
}

func NewAccountService() *AccountService {
	return &AccountService{
		accountRespository: respositories.NewAccountRespository(),
	}
}
