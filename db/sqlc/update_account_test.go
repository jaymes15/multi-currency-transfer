package db

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"lemfi/simplebank/util"
)

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	balance := decimal.NewFromFloat(util.RandomFloat(10.0, 1000.0)).Round(2)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: balance,
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
