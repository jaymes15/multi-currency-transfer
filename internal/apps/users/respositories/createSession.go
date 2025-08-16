package users

import (
	"context"
	"time"

	db "lemfi/simplebank/db/sqlc"

	"github.com/google/uuid"
)

func (userRespository *UserRespository) CreateSession(username string, refreshTokenID uuid.UUID, refreshToken string, expiresAt time.Time) error {
	arg := db.CreateSessionParams{
		ID:           refreshTokenID,
		Username:     username,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	_, err := userRespository.queries.CreateSession(context.Background(), arg)
	return err
}
