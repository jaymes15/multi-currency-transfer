package accounts

import (
	respositories "lemfi/simplebank/internal/apps/accounts/respositories"
)

type AccountService struct {
	accountRespository respositories.AccountRespositoryInterface
}

func NewAccountService(respository respositories.AccountRespositoryInterface) *AccountService {
	return &AccountService{
		accountRespository: respository,
	}
}
