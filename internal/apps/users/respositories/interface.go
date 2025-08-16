package users

import (
	"time"

	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/users/requests"

	"github.com/google/uuid"
)

type UserRespositoryInterface interface {
	CreateUser(payload requests.CreateUserRequest) (db.CreateUserRow, error)
	GetUserHashedPassword(username string) (string, error)
	GetUser(username string) (db.GetUserRow, error)
	CreateSession(username string, refreshTokenID uuid.UUID, refreshToken string, expiresAt time.Time) error
	GetSession(refreshTokenID uuid.UUID) (db.GetSessionRow, error)
	BlockSession(sessionID uuid.UUID) error
}
