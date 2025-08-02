package transfers

import (
	"context"

	dbConnection "lemfi/simplebank/db"
	db "lemfi/simplebank/db/sqlc"
)

type TransferRespository struct {
	context context.Context
	queries db.Store
}

func NewTransferRespository() *TransferRespository {
	return &TransferRespository{
		context: context.Background(),
		queries: db.NewStore(dbConnection.GetPostgresDBConnection()),
	}
}
