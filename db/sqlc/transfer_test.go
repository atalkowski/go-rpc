package db

import (
	"atalkowski/go-rpc/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccountID int64, toAccountID int64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	require.Equal(t, fromAccountID, transfer.FromAccountID)
	require.Equal(t, toAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	createRandomTransfer(t, fromAccount.ID, toAccount.ID)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	transfer := createRandomTransfer(t, fromAccount.ID, toAccount.ID)

	transfer1, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer1)
	require.NotZero(t, transfer1.ID)
	require.Equal(t, fromAccount.ID, transfer1.FromAccountID)
	require.Equal(t, toAccount.ID, transfer1.ToAccountID)
	require.Equal(t, transfer.Amount, transfer1.Amount)
	require.NotZero(t, transfer1.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fromAccount.ID, toAccount.ID)
	}
	arg := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotZero(t, transfer.ID)
		require.True(t, fromAccount.ID == transfer.FromAccountID || fromAccount.ID == transfer.ToAccountID)
		require.True(t, toAccount.ID == transfer.FromAccountID || toAccount.ID == transfer.ToAccountID)
		require.NotZero(t, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)
	}

}
