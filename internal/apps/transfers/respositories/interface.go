package transfers

import (
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/transfers/requests"
)

type TransferRespositoryInterface interface {
	MakeTransfer(payload requests.MakeTransferRequest) (db.TransferTxResult, error)
}
