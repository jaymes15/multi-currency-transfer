package db

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"lemfi/simplebank/util"
)

func TestUpdateExchangeRate(t *testing.T) {
	exchangeRate1 := createRandomExchangeRate(t)

	newRate := decimal.NewFromFloat(util.RandomFloat(0.1, 10.0)).Round(8)

	arg := UpdateExchangeRateParams{
		FromCurrency: exchangeRate1.FromCurrency,
		ToCurrency:   exchangeRate1.ToCurrency,
		Rate:         newRate,
	}

	exchangeRate2, err := testQueries.UpdateExchangeRate(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, exchangeRate2)

	require.Equal(t, exchangeRate1.ID, exchangeRate2.ID)
	require.Equal(t, exchangeRate1.FromCurrency, exchangeRate2.FromCurrency)
	require.Equal(t, exchangeRate1.ToCurrency, exchangeRate2.ToCurrency)
	require.Equal(t, newRate, exchangeRate2.Rate)
}
