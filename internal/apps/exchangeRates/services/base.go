package exchangeRates

import (
	respositories "lemfi/simplebank/internal/apps/exchangeRates/respositories"
)

type ExchangeRateService struct {
	exchangeRateRepository respositories.ExchangeRateRepositoryInterface
}

func NewExchangeRateService(exchangeRateRepository respositories.ExchangeRateRepositoryInterface) *ExchangeRateService {
	return &ExchangeRateService{
		exchangeRateRepository: exchangeRateRepository,
	}
}
