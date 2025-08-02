package db

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestAddAccountBalance(t *testing.T) {
	account := createRandomAccount(t)

	amount := decimal.NewFromFloat(10.0).Round(2)

	updatedBalance, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     account.ID,
		Amount: amount,
	})
	require.NoError(t, err)

	// Get the updated account to verify the balance
	updatedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, updatedBalance, updatedAccount.Balance)
}

func TestAddAccountBalanceNegative(t *testing.T) {
	account := createRandomAccount(t)

	amount := decimal.NewFromFloat(-10.0).Round(2)

	updatedBalance, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     account.ID,
		Amount: amount,
	})
	require.NoError(t, err)

	// Get the updated account to verify the balance
	updatedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, updatedBalance, updatedAccount.Balance)
}

func TestAddAccountBalanceZero(t *testing.T) {
	account := createRandomAccount(t)

	amount := decimal.NewFromFloat(0.0).Round(2)

	updatedBalance, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     account.ID,
		Amount: amount,
	})
	require.NoError(t, err)

	// Get the updated account to verify the balance
	updatedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, updatedBalance, updatedAccount.Balance) // Balance should remain the same
}

func TestAddAccountBalanceNotFound(t *testing.T) {
	amount := decimal.NewFromFloat(10.0).Round(2)

	_, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     99999, // Non-existent account
		Amount: amount,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rows in result set")
}
