package db

import (
	"context"
	"lemfi/simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		arg := CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomMoney(),
		}
		_, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}
