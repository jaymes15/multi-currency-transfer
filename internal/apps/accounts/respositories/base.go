package accounts

import (
	"context"

	dbConnection "lemfi/simplebank/db"
	db "lemfi/simplebank/db/sqlc"
)

type AccountRespository struct {
	context context.Context
	queries db.Store
}

func NewAccountRespository() *AccountRespository {
	return &AccountRespository{
		context: context.Background(),
		queries: db.NewStore(dbConnection.GetPostgresDBConnection()),
	}
}
