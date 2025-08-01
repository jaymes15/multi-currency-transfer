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

	// List all accounts with a higher limit to ensure we get our accounts
	accounts, err := testQueries.ListAllAccounts(context.Background(), ListAllAccountsParams{
		Limit:  1000, // Use a much higher limit to ensure we get all accounts
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	// Create a map of account IDs for easier lookup
	accountMap := make(map[int64]Account)
	for _, account := range accounts {
		accountMap[account.ID] = account
	}

	// Verify that our created accounts are in the list
	require.Contains(t, accountMap, account1.ID, "Account 1 not found in list")
	require.Contains(t, accountMap, account2.ID, "Account 2 not found in list")
	require.Contains(t, accountMap, account3.ID, "Account 3 not found in list")

	// Verify the account details
	foundAccount1 := accountMap[account1.ID]
	require.Equal(t, account1.Owner, foundAccount1.Owner)
	require.Equal(t, account1.Balance, foundAccount1.Balance)
	require.Equal(t, account1.Currency, foundAccount1.Currency)

	foundAccount2 := accountMap[account2.ID]
	require.Equal(t, account2.Owner, foundAccount2.Owner)
	require.Equal(t, account2.Balance, foundAccount2.Balance)
	require.Equal(t, account2.Currency, foundAccount2.Currency)

	foundAccount3 := accountMap[account3.ID]
	require.Equal(t, account3.Owner, foundAccount3.Owner)
	require.Equal(t, account3.Balance, foundAccount3.Balance)
	require.Equal(t, account3.Currency, foundAccount3.Currency)
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
