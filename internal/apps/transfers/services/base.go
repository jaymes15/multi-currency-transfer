package transfers

import (
	exchangeRateService "lemfi/simplebank/internal/apps/exchangeRates/services"
	respositories "lemfi/simplebank/internal/apps/transfers/respositories"
)

type TransferService struct {
	transferRespository respositories.TransferRespositoryInterface
	exchangeRateService *exchangeRateService.ExchangeRateService
}

func NewTransferService(
	respository respositories.TransferRespositoryInterface,
	exchangeRateService *exchangeRateService.ExchangeRateService,
) *TransferService {
	return &TransferService{
		transferRespository: respository,
		exchangeRateService: exchangeRateService,
	}
}
