package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddAccountBalance(t *testing.T) {
	account := createRandomAccount(t)

	amount := int64(10)
	updatedBalance, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     account.ID,
		Amount: amount,
	})
	require.NoError(t, err)

	expectedBalance := account.Balance + amount
	require.Equal(t, expectedBalance, updatedBalance)
}

func TestAddAccountBalanceNegative(t *testing.T) {
	account := createRandomAccount(t)

	amount := int64(-10)
	updatedBalance, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     account.ID,
		Amount: amount,
	})
	require.NoError(t, err)

	expectedBalance := account.Balance + amount
	require.Equal(t, expectedBalance, updatedBalance)
}

func TestAddAccountBalanceZero(t *testing.T) {
	account := createRandomAccount(t)

	amount := int64(0)
	updatedBalance, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     account.ID,
		Amount: amount,
	})
	require.NoError(t, err)

	expectedBalance := account.Balance + amount
	require.Equal(t, expectedBalance, updatedBalance) // Balance should remain the same
}

func TestAddAccountBalanceNotFound(t *testing.T) {
	_, err := testQueries.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID:     99999, // Non-existent account
		Amount: 10,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no rows in result set")
}
