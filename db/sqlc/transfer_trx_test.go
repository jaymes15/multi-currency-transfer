package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"lemfi/simplebank/util"
)

// Utility functions for decimal.Decimal arithmetic
func decimalToFloat64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	existed := make(map[int]bool)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Printf(">>>>> start transfer from account %d to account %d\n", account1.ID, account2.ID)
	fmt.Printf(">>>>> initial balance account1: %s, account2: %s\n", account1.Balance.String(), account2.Balance.String())

	// run n concurrent transfer transactions
	n := 5
	amount := decimal.NewFromInt(10).Round(2)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func(i int) {

			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID:   account1.ID,
				ToAccountID:     account2.ID,
				Amount:          amount,
				ConvertedAmount: amount, // Same currency, so converted amount equals original amount
				ExchangeRate:    decimal.NewFromInt(1).Round(8),
				FromCurrency:    account1.Currency,
				ToCurrency:      account2.Currency,
			})

			if err != nil {
				fmt.Printf(">>>>> goroutine %d: transfer failed: %v\n", i, err)
			} else {
				fmt.Printf(">>>>> goroutine %d: transfer completed successfully\n", i)
			}

			errs <- err
			results <- result
		}(i)
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		fmt.Printf(">>>>> checking result %d\n", i+1)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, amount.Neg(), fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check to entries
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check from account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		// check to account
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check from account balance
		diff1 := account1.Balance.Sub(fromAccount.Balance)
		require.True(t, diff1.GreaterThan(decimal.Zero))
		require.True(t, diff1.Mod(amount).IsZero()) // 1 * amount, 2 * amount, 3 * amount, ...

		// check to account balance
		diff2 := toAccount.Balance.Sub(account2.Balance)
		require.True(t, diff2.GreaterThan(decimal.Zero))
		require.True(t, diff2.Mod(amount).IsZero())

		require.Equal(t, diff1, diff2)

		k := int(diff1.Div(amount).IntPart())
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
		fmt.Printf(">>>>> transfer %d: k=%d, diff1=%s, diff2=%s\n", i+1, k, diff1.String(), diff2.String())
	}

	// check the final balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">>>>> final balance account1: %s, account2: %s\n", updatedAccount1.Balance.String(), updatedAccount2.Balance.String())

	// Verify final balances
	expectedBalance1 := account1.Balance.Sub(amount.Mul(decimal.NewFromInt(int64(n))))
	expectedBalance2 := account2.Balance.Add(amount.Mul(decimal.NewFromInt(int64(n))))

	require.Equal(t, expectedBalance1, updatedAccount1.Balance)
	require.Equal(t, expectedBalance2, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Printf(">>>>> start transfer from account %d to account %d\n", account1.ID, account2.ID)
	fmt.Printf(">>>>> initial balance account1: %s, account2: %s\n", account1.Balance.String(), account2.Balance.String())

	// run n concurrent transfer transactions
	n := 10
	amount := decimal.NewFromInt(10).Round(2)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func(i int) {
			// Get the actual currencies for the accounts being used
			fromAccount, err := testQueries.GetAccount(context.Background(), fromAccountID)
			if err != nil {
				errs <- err
				return
			}
			toAccount, err := testQueries.GetAccount(context.Background(), toAccountID)
			if err != nil {
				errs <- err
				return
			}

			_, err = store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID:   fromAccountID,
				ToAccountID:     toAccountID,
				Amount:          amount,
				ConvertedAmount: amount, // Same currency, so converted amount equals original amount
				ExchangeRate:    decimal.NewFromInt(1).Round(8),
				FromCurrency:    fromAccount.Currency,
				ToCurrency:      toAccount.Currency,
			})

			if err != nil {
				fmt.Printf(">>>>> goroutine %d: transfer failed: %v\n", i, err)
			} else {
				fmt.Printf(">>>>> goroutine %d: transfer completed successfully\n", i)
			}

			errs <- err
		}(i)
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		fmt.Printf(">>>>> checking result %d\n", i+1)

	}

	// check the final balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">>>>> final balance account1: %s, account2: %s\n", updatedAccount1.Balance.String(), updatedAccount2.Balance.String())

	// Verify final balances (should be the same since transfers are bidirectional)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}

func TestTransferTxMultiCurrency(t *testing.T) {
	store := NewStore(testDB)

	// Create accounts with different currencies
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Ensure accounts have different currencies
	for account1.Currency == account2.Currency {
		account2 = createRandomAccount(t)
	}

	// Try to get existing exchange rate, or create a new one
	var exchangeRate ExchangeRate
	existingRate, err := testQueries.GetExchangeRate(context.Background(), GetExchangeRateParams{
		FromCurrency: account1.Currency,
		ToCurrency:   account2.Currency,
	})

	if err != nil {
		// Exchange rate doesn't exist, create a new one
		rate := decimal.NewFromFloat(util.RandomFloat(0.1, 10.0))

		exchangeRateArg := CreateExchangeRateParams{
			FromCurrency: account1.Currency,
			ToCurrency:   account2.Currency,
			Rate:         rate,
		}

		exchangeRate, err = testQueries.CreateExchangeRate(context.Background(), exchangeRateArg)
		require.NoError(t, err)
	} else {
		exchangeRate = existingRate
	}

	require.NotEmpty(t, exchangeRate)

	fmt.Printf(">>>>> start multi-currency transfer from account %d (%s) to account %d (%s)\n",
		account1.ID, account1.Currency, account2.ID, account2.Currency)
	fmt.Printf(">>>>> initial balance account1: %s %s, account2: %s %s\n",
		account1.Balance.String(), account1.Currency, account2.Balance.String(), account2.Currency)
	fmt.Printf(">>>>> exchange rate: 1 %s = %s %s\n",
		account1.Currency, exchangeRate.Rate.String(), account2.Currency)

	// Use a smaller amount to ensure it's less than the account balance
	amount := decimal.NewFromInt(10).Round(2)

	// Calculate converted amount and round to 2 decimal places
	convertedAmount := amount.Mul(exchangeRate.Rate).Round(2)

	// Perform the multi-currency transfer
	result, err := store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID:   account1.ID,
		ToAccountID:     account2.ID,
		Amount:          amount,
		ConvertedAmount: convertedAmount,
		ExchangeRate:    exchangeRate.Rate,
		FromCurrency:    account1.Currency,
		ToCurrency:      account2.Currency,
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)

	fmt.Printf(">>>>> transfer completed successfully\n")

	// Check transfer
	transfer := result.Transfer
	require.NotEmpty(t, transfer)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)
	require.NotZero(t, transfer.ID)

	// Verify the transfer was created in the database
	createdTransfer, err := store.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.Equal(t, transfer.ID, createdTransfer.ID)

	// Check from entry
	fromEntry := result.FromEntry
	require.NotEmpty(t, fromEntry)
	require.Equal(t, account1.ID, fromEntry.AccountID)
	require.Equal(t, amount.Neg(), fromEntry.Amount)
	require.NotZero(t, fromEntry.ID)

	// Verify the from entry was created in the database
	createdFromEntry, err := store.GetEntry(context.Background(), fromEntry.ID)
	require.NoError(t, err)
	require.Equal(t, fromEntry.ID, createdFromEntry.ID)

	// Check to entry
	toEntry := result.ToEntry
	require.NotEmpty(t, toEntry)
	require.Equal(t, account2.ID, toEntry.AccountID)
	require.Equal(t, convertedAmount, toEntry.Amount) // Should be converted amount
	require.NotZero(t, toEntry.ID)

	// Verify the to entry was created in the database
	createdToEntry, err := store.GetEntry(context.Background(), toEntry.ID)
	require.NoError(t, err)
	require.Equal(t, toEntry.ID, createdToEntry.ID)

	// Check from account
	fromAccount := result.FromAccount
	require.NotEmpty(t, fromAccount)
	require.Equal(t, account1.ID, fromAccount.ID)

	// Check to account
	toAccount := result.ToAccount
	require.NotEmpty(t, toAccount)
	require.Equal(t, account2.ID, toAccount.ID)

	// Check from account balance (should be reduced by the original amount)
	diff1 := account1.Balance.Sub(fromAccount.Balance)
	require.Equal(t, amount, diff1)

	// Check to account balance (should be increased by the converted amount)
	diff2 := toAccount.Balance.Sub(account2.Balance)
	require.Equal(t, convertedAmount, diff2)

	fmt.Printf(">>>>> balance change account1: -%s %s, account2: +%s %s\n",
		diff1.String(), account1.Currency, diff2.String(), account2.Currency)

	// Verify final balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">>>>> final balance account1: %s %s, account2: %s %s\n",
		updatedAccount1.Balance.String(), updatedAccount1.Currency, updatedAccount2.Balance.String(), updatedAccount2.Currency)

	// Check final balances with proper arithmetic
	expectedBalance1 := account1.Balance.Sub(amount)
	expectedBalance2 := account2.Balance.Add(convertedAmount)

	require.Equal(t, expectedBalance1, updatedAccount1.Balance)
	require.Equal(t, expectedBalance2, updatedAccount2.Balance)

	// Verify currencies haven't changed
	require.Equal(t, account1.Currency, updatedAccount1.Currency)
	require.Equal(t, account2.Currency, updatedAccount2.Currency)
}
