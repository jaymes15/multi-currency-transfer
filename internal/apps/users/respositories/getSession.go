package users

import (
	db "lemfi/simplebank/db/sqlc"

	"github.com/google/uuid"
)

func (r *UserRespository) GetSession(refreshTokenID uuid.UUID) (db.GetSessionRow, error) {
	return r.queries.GetSession(r.context, refreshTokenID)
}
