package transfers

import "lemfi/simplebank/internal/apps/core"

// Predefined client errors
var (
	ErrSameAccountTransfer = core.ClientError{
		Message: "cannot transfer to the same account",
		Status:  400,
	}
	ErrInvalidAmount = core.ClientError{
		Message: "transfer amount must be positive",
		Status:  400,
	}
	ErrFromAccountNotFound = core.ClientError{
		Message: "from account not found",
		Status:  400,
	}
	ErrToAccountNotFound = core.ClientError{
		Message: "to account not found",
		Status:  400,
	}
	ErrFromAccountCurrencyMismatch = core.ClientError{
		Message: "from account currency mismatch",
		Status:  400,
	}
	ErrToAccountCurrencyMismatch = core.ClientError{
		Message: "to account currency mismatch",
		Status:  400,
	}
	ErrInsufficientBalance = core.ClientError{
		Message: "insufficient balance",
		Status:  400,
	}
)
