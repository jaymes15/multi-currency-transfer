package transfers

import (
	requests "lemfi/simplebank/internal/apps/transfers/requests"
	responses "lemfi/simplebank/internal/apps/transfers/responses"
)

type TransferServiceInterface interface {
	MakeTransfer(payload requests.MakeTransferRequest) (responses.MakeTransferResponse, error)
}
