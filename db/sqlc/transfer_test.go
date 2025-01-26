package db

import (
	"context"
	"testing"
	"time"

	"github.com/reynaldineo/simplebank-go/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, accountFrom, accountTo Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        utils.RandomMoney(),
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

func TestCreateTransfer(t *testing.T) {
	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)

	createRandomTransfer(t, accountFrom, accountTo)
}

func TestGeTransfer(t *testing.T) {
	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)

	transfer1 := createRandomTransfer(t, accountFrom, accountTo)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, accountFrom, accountTo)
	}

	arg := ListTransfersParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == accountFrom.ID)
		require.True(t, transfer.ToAccountID == accountTo.ID)
	}

}
