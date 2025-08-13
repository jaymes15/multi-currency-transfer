package users

import (
	"context"

	dbConnection "lemfi/simplebank/db"
	db "lemfi/simplebank/db/sqlc"
)

type UserRespository struct {
	context context.Context
	queries db.Store
}

func NewUserRespository() *UserRespository {
	return &UserRespository{
		context: context.Background(),
		queries: db.NewStore(dbConnection.GetPostgresDBConnection()),
	}
}
