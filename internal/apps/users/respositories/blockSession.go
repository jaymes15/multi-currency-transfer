package users

import (
	"github.com/google/uuid"
)

func (r *UserRespository) BlockSession(sessionID uuid.UUID) error {
	return r.queries.BlockSession(r.context, sessionID)
}
