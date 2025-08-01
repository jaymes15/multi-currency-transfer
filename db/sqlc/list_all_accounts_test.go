package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListAllAccounts(t *testing.T) {
	// Create multiple accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	account3 := createRandomAccount(t)

	// List all accounts
	accounts, err := testQueries.ListAllAccounts(context.Background(), ListAllAccountsParams{
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	// Verify that our created accounts are in the list
	found1, found2, found3 := false, false, false
	for _, account := range accounts {
		if account.ID == account1.ID {
			found1 = true
			require.Equal(t, account1.Owner, account.Owner)
			require.Equal(t, account1.Balance, account.Balance)
			require.Equal(t, account1.Currency, account.Currency)
		}
		if account.ID == account2.ID {
			found2 = true
			require.Equal(t, account2.Owner, account.Owner)
			require.Equal(t, account2.Balance, account.Balance)
			require.Equal(t, account2.Currency, account.Currency)
		}
		if account.ID == account3.ID {
			found3 = true
			require.Equal(t, account3.Owner, account.Owner)
			require.Equal(t, account3.Balance, account.Balance)
			require.Equal(t, account3.Currency, account.Currency)
		}
	}

	require.True(t, found1, "Account 1 not found in list")
	require.True(t, found2, "Account 2 not found in list")
	require.True(t, found3, "Account 3 not found in list")
}

func TestListAllAccountsWithLimit(t *testing.T) {
	// Create multiple accounts to ensure we have data
	createRandomAccount(t)
	createRandomAccount(t)
	createRandomAccount(t)

	// List accounts with limit 2
	accounts, err := testQueries.ListAllAccounts(context.Background(), ListAllAccountsParams{
		Limit:  2,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, accounts, 2)

	// Verify accounts are ordered by ID
	require.True(t, accounts[0].ID <= accounts[1].ID, "Accounts should be ordered by ID")
}

func TestListAllAccountsWithOffset(t *testing.T) {
	// Create multiple accounts to ensure we have data
	firstAccount := createRandomAccount(t)
	createRandomAccount(t)
	createRandomAccount(t)

	// List accounts with offset 1
	accounts, err := testQueries.ListAllAccounts(context.Background(), ListAllAccountsParams{
		Limit:  10,
		Offset: 1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	// Verify that the first account in the result is not the first created account
	// (due to offset)
	require.NotEqual(t, firstAccount.ID, accounts[0].ID, "First account should be skipped due to offset")
}

func TestListAllAccountsEmpty(t *testing.T) {
	// This test assumes the database might be empty or we can't guarantee it's empty
	// In a real scenario, you might want to clean the database before this test
	accounts, err := testQueries.ListAllAccounts(context.Background(), ListAllAccountsParams{
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)
	// The list might be empty or contain data, both are valid
	require.NotNil(t, accounts)
}
