package users

import (
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/users/requests"
)

type UserRespositoryInterface interface {
	CreateUser(payload requests.CreateUserRequest) (db.CreateUserRow, error)
}
