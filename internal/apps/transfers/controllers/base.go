package transfers

import (
	services "lemfi/simplebank/internal/apps/transfers/services"
)

type TransferController struct {
	transferService services.TransferServiceInterface
}

func NewTransferController(service services.TransferServiceInterface) *TransferController {
	return &TransferController{
		transferService: service,
	}
}
