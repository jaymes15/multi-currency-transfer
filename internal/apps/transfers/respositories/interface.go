package transfers

import (
	db "lemfi/simplebank/db/sqlc"
	requests "lemfi/simplebank/internal/apps/transfers/requests"

	"github.com/shopspring/decimal"
)

type TransferRespositoryInterface interface {
	MakeTransfer(payload requests.MakeTransferRequest, convertedAmount decimal.Decimal, exchangeRate decimal.Decimal) (db.TransferTxResult, error)
}
