package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	store := NewStore(testDB)

	accountTo := createRandomAccount(t)
	accountFrom := createRandomAccount(t)
	fmt.Println(">> Before:", accountFrom.Balance, accountTo.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: accountFrom.ID,
				ToAccountID:   accountTo.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.Equal(t, accountFrom.ID, transfer.FromAccountID)
		require.Equal(t, accountTo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntries := result.FromEntry
		require.NotEmpty(t, fromEntries)
		require.Equal(t, accountFrom.ID, fromEntries.AccountID)
		require.Equal(t, -amount, fromEntries.Amount)
		require.NotZero(t, fromEntries.ID)
		require.NotZero(t, fromEntries.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntries.ID)
		require.NoError(t, err)

		toEntries := result.ToEntry
		require.NotEmpty(t, toEntries)
		require.Equal(t, accountTo.ID, toEntries.AccountID)
		require.Equal(t, amount, toEntries.Amount)
		require.NotZero(t, toEntries.ID)
		require.NotZero(t, toEntries.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntries.ID)
		require.NoError(t, err)

		// check account
		fromAcount := result.FromAccount
		require.NotEmpty(t, fromAcount)
		require.Equal(t, accountFrom.ID, fromAcount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountTo.ID, toAccount.ID)

		// check account balance
		fmt.Println(">> Transferred:", fromAcount.Balance, toAccount.Balance)
		diff1 := accountFrom.Balance - fromAcount.Balance // amount transferred
		diff2 := toAccount.Balance - accountTo.Balance    // amount received
		require.Equal(t, diff1, diff2)                    // to check amount transferred and received are equal
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updateAccountFrom, err := testQueries.GetAccount(context.Background(), accountFrom.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccountFrom)

	updateAccountTo, err := testQueries.GetAccount(context.Background(), accountTo.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccountTo)

	fmt.Println(">> After:", updateAccountFrom.Balance, updateAccountTo.Balance)
	require.Equal(t, accountFrom.Balance-int64(n)*amount, updateAccountFrom.Balance)
	require.Equal(t, accountTo.Balance+int64(n)*amount, updateAccountTo.Balance)
}
