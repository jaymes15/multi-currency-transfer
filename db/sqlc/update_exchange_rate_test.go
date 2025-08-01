package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	"lemfi/simplebank/util"
)

func TestUpdateExchangeRate(t *testing.T) {
	exchangeRate1 := createRandomExchangeRate(t)

	newRate := util.RandomFloat(0.1, 5000.0)

	var numericRate pgtype.Numeric
	numericRate.Scan(fmt.Sprintf("%.8f", newRate))

	arg := UpdateExchangeRateParams{
		FromCurrency: exchangeRate1.FromCurrency,
		ToCurrency:   exchangeRate1.ToCurrency,
		Rate:         numericRate,
	}

	exchangeRate2, err := testQueries.UpdateExchangeRate(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, exchangeRate2)

	require.Equal(t, exchangeRate1.ID, exchangeRate2.ID)
	require.Equal(t, exchangeRate1.FromCurrency, exchangeRate2.FromCurrency)
	require.Equal(t, exchangeRate1.ToCurrency, exchangeRate2.ToCurrency)
	require.InDelta(t, newRate, exchangeRate2.Rate, 0.000001)
}
