package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetExchangeRate(t *testing.T) {
	exchangeRate1 := createRandomExchangeRate(t)
	exchangeRate2, err := testQueries.GetExchangeRate(context.Background(), GetExchangeRateParams{
		FromCurrency: exchangeRate1.FromCurrency,
		ToCurrency:   exchangeRate1.ToCurrency,
	})
	require.NoError(t, err)
	require.NotEmpty(t, exchangeRate2)

	require.Equal(t, exchangeRate1.ID, exchangeRate2.ID)
	require.Equal(t, exchangeRate1.FromCurrency, exchangeRate2.FromCurrency)
	require.Equal(t, exchangeRate1.ToCurrency, exchangeRate2.ToCurrency)
	require.Equal(t, exchangeRate1.Rate, exchangeRate2.Rate)
}
