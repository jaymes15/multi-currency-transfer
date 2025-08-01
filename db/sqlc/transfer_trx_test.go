package db

import (
	"context"
	"fmt"
	"testing"

	"lemfi/simplebank/util"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	existed := make(map[int]bool)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Printf(">>>>> start transfer from account %d to account %d\n", account1.ID, account2.ID)
	fmt.Printf(">>>>> initial balance account1: %d, account2: %d\n", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func(i int) {

			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
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
		require.Equal(t, -amount, fromEntry.Amount)
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
		diff1 := account1.Balance - fromAccount.Balance
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ...

		// check to account balance
		diff2 := toAccount.Balance - account2.Balance
		require.True(t, diff2 > 0)
		require.True(t, diff2%amount == 0)

		require.Equal(t, diff1, diff2)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
		fmt.Printf(">>>>> transfer %d: k=%d, diff1=%d, diff2=%d\n", i+1, k, diff1, diff2)
	}

	// check the final balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">>>>> final balance account1: %d, account2: %d\n", updatedAccount1.Balance, updatedAccount2.Balance)
	fmt.Printf(">>>>> expected balance account1: %d, account2: %d\n", account1.Balance-int64(n)*amount, account2.Balance+int64(n)*amount)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Printf(">>>>> start transfer from account %d to account %d\n", account1.ID, account2.ID)
	fmt.Printf(">>>>> initial balance account1: %d, account2: %d\n", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func(i int) {

			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
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

	fmt.Printf(">>>>> final balance account1: %d, account2: %d\n", updatedAccount1.Balance, updatedAccount2.Balance)
	fmt.Printf(">>>>> expected balance account1: %d, account2: %d\n", account1.Balance-int64(n)*amount, account2.Balance+int64(n)*amount)

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
	var exchangeRate GetExchangeRateRow
	existingRate, err := testQueries.GetExchangeRate(context.Background(), GetExchangeRateParams{
		FromCurrency: account1.Currency,
		ToCurrency:   account2.Currency,
	})

	if err != nil {
		// Exchange rate doesn't exist, create a new one
		rate := util.RandomFloat(0.1, 5000.0)
		exchangeRateArg := CreateExchangeRateParams{
			FromCurrency: account1.Currency,
			ToCurrency:   account2.Currency,
			Rate:         fmt.Sprintf("%.8f", rate),
		}

		createdRate, err := testQueries.CreateExchangeRate(context.Background(), exchangeRateArg)
		require.NoError(t, err)

		// Convert CreateExchangeRateRow to GetExchangeRateRow
		exchangeRate = GetExchangeRateRow{
			ID:           createdRate.ID,
			FromCurrency: createdRate.FromCurrency,
			ToCurrency:   createdRate.ToCurrency,
			Rate:         createdRate.Rate,
			CreatedAt:    createdRate.CreatedAt,
		}
	} else {
		exchangeRate = existingRate
	}

	require.NotEmpty(t, exchangeRate)

	fmt.Printf(">>>>> start multi-currency transfer from account %d (%s) to account %d (%s)\n",
		account1.ID, account1.Currency, account2.ID, account2.Currency)
	fmt.Printf(">>>>> initial balance account1: %d %s, account2: %d %s\n",
		account1.Balance, account1.Currency, account2.Balance, account2.Currency)
	fmt.Printf(">>>>> exchange rate: 1 %s = %.8f %s\n",
		account1.Currency, exchangeRate.Rate, account2.Currency)

	amount := int64(100)

	// Perform the multi-currency transfer
	result, err := store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
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
	require.Equal(t, -amount, fromEntry.Amount)
	require.NotZero(t, fromEntry.ID)

	// Verify the from entry was created in the database
	createdFromEntry, err := store.GetEntry(context.Background(), fromEntry.ID)
	require.NoError(t, err)
	require.Equal(t, fromEntry.ID, createdFromEntry.ID)

	// Check to entry
	toEntry := result.ToEntry
	require.NotEmpty(t, toEntry)
	require.Equal(t, account2.ID, toEntry.AccountID)
	require.Equal(t, amount, toEntry.Amount)
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
	diff1 := account1.Balance - fromAccount.Balance
	require.Equal(t, amount, diff1)

	// Check to account balance (should be increased by the original amount)
	diff2 := toAccount.Balance - account2.Balance
	require.Equal(t, amount, diff2)

	fmt.Printf(">>>>> balance change account1: -%d %s, account2: +%d %s\n",
		diff1, account1.Currency, diff2, account2.Currency)

	// Verify final balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf(">>>>> final balance account1: %d %s, account2: %d %s\n",
		updatedAccount1.Balance, updatedAccount1.Currency, updatedAccount2.Balance, updatedAccount2.Currency)

	require.Equal(t, account1.Balance-amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+amount, updatedAccount2.Balance)

	// Verify currencies haven't changed
	require.Equal(t, account1.Currency, updatedAccount1.Currency)
	require.Equal(t, account2.Currency, updatedAccount2.Currency)
}
