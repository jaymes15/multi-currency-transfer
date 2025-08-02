package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"lemfi/simplebank/util"
)

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Create transfers from account1 to account2
	for i := 0; i < 5; i++ {
		amount := decimal.NewFromFloat(util.RandomFloat(10.0, 1000.0)).Round(2)

		arg := CreateTransferParams{
			FromAccountID:   account1.ID,
			ToAccountID:     account2.ID,
			Amount:          amount,
			ConvertedAmount: amount, // For simple transfers, use same amount
			ExchangeRate:    decimal.NewFromFloat(1.0).Round(8),
			FromCurrency:    pgtype.Text{},
			ToCurrency:      pgtype.Text{},
		}
		_, err := testQueries.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	// Create transfers from account2 to account1
	for i := 0; i < 3; i++ {
		amount := decimal.NewFromFloat(util.RandomFloat(10.0, 1000.0)).Round(2)

		arg := CreateTransferParams{
			FromAccountID:   account2.ID,
			ToAccountID:     account1.ID,
			Amount:          amount,
			ConvertedAmount: amount, // For simple transfers, use same amount
			ExchangeRate:    decimal.NewFromFloat(1.0).Round(8),
			FromCurrency:    pgtype.Text{},
			ToCurrency:      pgtype.Text{},
		}
		_, err := testQueries.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	// List transfers for account1 (both sent and received)
	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         10,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	// Verify all transfers involve account1
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}
