package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func checkValidEntry(t *testing.T, store *Store, entry Entry, accountId int64, amount int64) {
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.Equal(t, accountId, entry.AccountID)
	require.Equal(t, amount, entry.Amount)
	require.NotZero(t, entry.CreatedAt)
	_, err := store.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
}

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check the entries and FromAccount was debited and ToAccount credited with Anount
		checkValidEntry(t, store, result.FromEntry, account1.ID, -amount)
		checkValidEntry(t, store, result.ToEntry, account2.ID, amount)

		// TODO : check Accounts were correctly debited and credited .. when we have built that!!
	}
}
