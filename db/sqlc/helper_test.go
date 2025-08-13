package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"lemfi/simplebank/util"
)

func createRandomUser(t *testing.T) CreateUserRow {
	// Generate random user data
	username := util.RandomOwner()                            // Reuse existing random owner function
	fullName := util.RandomOwner() + " " + util.RandomOwner() // Create a full name
	email := username + "@example.com"
	hashedPassword := "hashedpassword123" // Fixed password for testing

	arg := CreateUserParams{
		Username:       username,
		HashedPassword: hashedPassword,
		FullName:       fullName,
		Email:          email,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)

	return user
}

func createRandomAccount(t *testing.T) Account {
	// Create a user first since accounts now have foreign key constraint
	user := createRandomUser(t)

	// Create balance with 2 decimal places to match database precision
	balanceFloat := util.RandomFloat(10.0, 1000.0)
	balance := decimal.NewFromFloat(balanceFloat).Round(2)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  balance,
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

func createAccountWithCurrency(t *testing.T, currency string) Account {
	// Create a user first since accounts now have foreign key constraint
	user := createRandomUser(t)

	// Create balance with 2 decimal places to match database precision
	balanceFloat := util.RandomFloat(10.0, 1000.0)
	balance := decimal.NewFromFloat(balanceFloat).Round(2)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  balance,
		Currency: currency,
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

	// Create amount with 2 decimal places to match database precision
	amountFloat := util.RandomFloat(10.0, 1000.0)
	amount := decimal.NewFromFloat(amountFloat).Round(2)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    amount,
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

func createRandomTransfer(t *testing.T) Transfer {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	// Create amount with 2 decimal places to match database precision
	amountFloat := util.RandomFloat(10.0, 1000.0)
	amount := decimal.NewFromFloat(amountFloat).Round(2)

	// Create fee with 2 decimal places
	feeFloat := util.RandomFloat(1.0, 5.0)
	fee := decimal.NewFromFloat(feeFloat).Round(2)

	arg := CreateTransferParams{
		FromAccountID:   fromAccount.ID,
		ToAccountID:     toAccount.ID,
		Amount:          amount,
		ConvertedAmount: amount,                             // For simple transfers, use same amount
		ExchangeRate:    decimal.NewFromFloat(1.0).Round(8), // 8 decimal places for exchange rates
		FromCurrency:    pgtype.Text{String: fromAccount.Currency, Valid: true},
		ToCurrency:      pgtype.Text{String: toAccount.Currency, Valid: true},
		Fee:             fee,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, arg.Fee, transfer.Fee)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	// Convert CreateTransferRow to Transfer
	result := Transfer{
		ID:              transfer.ID,
		FromAccountID:   transfer.FromAccountID,
		ToAccountID:     transfer.ToAccountID,
		Amount:          transfer.Amount,
		ConvertedAmount: transfer.ConvertedAmount,
		ExchangeRate:    transfer.ExchangeRate,
		FromCurrency:    transfer.FromCurrency,
		ToCurrency:      transfer.ToCurrency,
		Fee:             transfer.Fee,
		CreatedAt:       transfer.CreatedAt,
	}

	return result
}

func createRandomExchangeRate(t *testing.T) ExchangeRate {
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
		return existingRate
	}

	// Create new exchange rate with 8 decimal places to match database precision
	rateFloat := util.RandomFloat(0.1, 10.0)
	rate := decimal.NewFromFloat(rateFloat).Round(8)

	arg := CreateExchangeRateParams{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
	}

	exchangeRate, err := testQueries.CreateExchangeRate(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, exchangeRate)

	require.Equal(t, arg.FromCurrency, exchangeRate.FromCurrency)
	require.Equal(t, arg.ToCurrency, exchangeRate.ToCurrency)
	require.Equal(t, arg.Rate, exchangeRate.Rate)
	require.NotZero(t, exchangeRate.ID)
	require.NotZero(t, exchangeRate.CreatedAt)

	return exchangeRate
}
