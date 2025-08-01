package accounts

import (
	"context"
	"database/sql"

	dbConnection "lemfi/simplebank/db"
	db "lemfi/simplebank/db/sqlc"
)

type AccountRespository struct {
	PostgresDB *sql.DB
	context    context.Context
	queries    *db.Queries
}

func NewAccountRespository() *AccountRespository {
	return &AccountRespository{
		PostgresDB: dbConnection.GetPostgresDBConnection(),
		context:    context.Background(),
		queries:    db.New(dbConnection.GetPostgresDBConnection()),
	}
}
