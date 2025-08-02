package transfers

import (
	respositories "lemfi/simplebank/internal/apps/transfers/respositories"
)

type TransferService struct {
	transferRespository respositories.TransferRespositoryInterface
}

func NewTransferService(respository respositories.TransferRespositoryInterface) *TransferService {
	return &TransferService{
		transferRespository: respository,
	}
}
