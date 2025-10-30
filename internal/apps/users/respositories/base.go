package users

import (
	"context"

	dbConnection "lemfi/simplebank/db"
	db "lemfi/simplebank/db/sqlc"

	"github.com/google/uuid"
)

// dbQuerier captures only the DB methods this repository needs.
type dbQuerier interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
	GetUserHashedPassword(ctx context.Context, username string) (string, error)
	GetUser(ctx context.Context, username string) (db.GetUserRow, error)
	CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (db.GetSessionRow, error)
	BlockSession(ctx context.Context, id uuid.UUID) error
}

type UserRespository struct {
	context context.Context
	queries dbQuerier
}

func NewUserRespository() *UserRespository {
	return &UserRespository{
		context: context.Background(),
		queries: db.New(dbConnection.GetPostgresDBConnection()),
	}
}
