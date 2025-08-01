package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	"lemfi/simplebank/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func createRandomTransfer(t *testing.T) CreateTransferRow {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomMoney(),
		ExchangeRate:  pgtype.Numeric{},
		FromCurrency:  pgtype.Text{},
		ToCurrency:    pgtype.Text{},
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func createRandomExchangeRate(t *testing.T) CreateExchangeRateRow {
	currencies := []string{"GBP", "NGN", "USD", "EUR"}
	fromCurrency := currencies[util.RandomInt(0, int64(len(currencies)-1))]
	toCurrency := currencies[util.RandomInt(0, int64(len(currencies)-1))]

	// Ensure different currencies
	for toCurrency == fromCurrency {
		toCurrency = currencies[util.RandomInt(0, int64(len(currencies)-1))]
	}

	// Try to get existing exchange rate first
	existingRate, err := testQueries.GetExchangeRate(context.Background(), GetExchangeRateParams{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	})

	if err == nil {
		// Exchange rate already exists, return it
		return CreateExchangeRateRow{
			ID:           existingRate.ID,
			FromCurrency: existingRate.FromCurrency,
			ToCurrency:   existingRate.ToCurrency,
			Rate:         existingRate.Rate,
			CreatedAt:    existingRate.CreatedAt,
		}
	}

	// Create new exchange rate
	rate := util.RandomFloat(0.1, 5000.0)

	var numericRate pgtype.Numeric
	numericRate.Scan(fmt.Sprintf("%.8f", rate))

	arg := CreateExchangeRateParams{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         numericRate,
	}

	exchangeRate, err := testQueries.CreateExchangeRate(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, exchangeRate)

	require.Equal(t, arg.FromCurrency, exchangeRate.FromCurrency)
	require.Equal(t, arg.ToCurrency, exchangeRate.ToCurrency)
	require.InDelta(t, rate, exchangeRate.Rate, 0.000001) // Use InDelta for float comparison
	require.NotZero(t, exchangeRate.ID)
	require.NotZero(t, exchangeRate.CreatedAt)

	return exchangeRate
}
