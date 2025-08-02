package accounts

import "lemfi/simplebank/internal/apps/core"

// Predefined client errors
var (
	ErrDuplicateAccount = core.ClientError{
		Message: "account already exists for this owner and currency",
		Status:  400,
	}
	ErrAccountNotFound = core.ClientError{
		Message: "account not found",
		Status:  404,
	}
)
