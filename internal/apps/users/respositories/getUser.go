package users

import (
	db "lemfi/simplebank/db/sqlc"
)

func (r *UserRespository) GetUser(username string) (db.GetUserRow, error) {
	return r.queries.GetUser(r.context, username)
}
