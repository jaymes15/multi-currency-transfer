package users

import "lemfi/simplebank/internal/apps/core"

// Predefined client errors
var (
	ErrDuplicateUsername = core.ClientError{
		Message: "username already exists",
		Status:  400,
	}
	ErrDuplicateEmail = core.ClientError{
		Message: "email already exists",
		Status:  400,
	}
	ErrUserNotFound = core.ClientError{
		Message: "user not found",
		Status:  404,
	}
)
