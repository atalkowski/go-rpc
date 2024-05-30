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
	fmt.Printf(">> before: from=(%v, %v cents) to=(%v, %v cents)\n", account1.ID, account1.Balance,
		account2.ID, account2.Balance)

	n := 4
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

	existed := make(map[int]bool)
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

		// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Check Account balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // amounts should be in increments of amount (between x1 and x5)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	fmt.Printf(">> after: from=(%v, %v cents) to=(%v, %v cents)\n", updatedAccount1.ID, updatedAccount1.Balance,
		updatedAccount2.ID, updatedAccount2.Balance)
}

/*
The following test will expose the deadlock error like this:
>> before: from=(392, 789 cents) to=(393, 188 cents)
--- FAIL: TestTransferDeadlockTx (1.05s)

	/Users/andy/wspaces/go/go-rpc/db/sqlc/store_test.go:136:
	    	Error Trace:	/Users/andy/wspaces/go/go-rpc/db/sqlc/store_test.go:136
	    	Error:      	Received unexpected error:
	    	            	pq: deadlock detected
	    	Test:       	TestTransferDeadlockTx

FAIL
FAIL	atalkowski/go-rpc/db/sqlc	1.439s
FAIL
TODO: fix this as discussed in the store.go
*/
func TestTransferDeadlockTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Printf(">> before: from=(%v, %v cents) to=(%v, %v cents)\n", account1.ID, account1.Balance,
		account2.ID, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		accountID1 := account1.ID
		accountID2 := account2.ID
		if i%2 == 1 {
			accountID1 = account2.ID
			accountID2 = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: accountID1,
				ToAccountID:   accountID2,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

	fmt.Printf(">> after: from=(%v, %v cents) to=(%v, %v cents)\n", updatedAccount1.ID, updatedAccount1.Balance,
		updatedAccount2.ID, updatedAccount2.Balance)
}
