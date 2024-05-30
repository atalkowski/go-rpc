package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries abd transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx (not exported!) exacutes a db function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // Use deft options set in DB rather than &sql.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx) // Pass the transaction rather than the Db
	err = fn(q)  // Invoke the intended sql
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb errL %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains input for transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer,omitempty"`
	FromAccount Account  `json:"from_account,omitempty"` // Updated FromAccount
	ToAccount   Account  `json:"to_account,omitempty"`   // Updated ToAccount
	FromEntry   Entry    `json:"from_entry,omitempty"`   // Recor dof money moving out
	ToEntry     Entry    `json:"to_entry,omitempty"`     // Record of money comint in
}

// var txKey = struct{}{} // Used with ctx to add a vble to the context for debugging

// Transfer performs all necessary updates to execute a money transfer between 2 accounts.
// It creates a transfer record, adds account entries and updates the accont balances for each account
// Within a single transcation
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// Note use of in line function here:-
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// txName := ctx.Value(txKey)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update the FromEntry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update the ToEntry
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = txnAvoidDeadlock(ctx, q, arg.FromAccountID, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = txnAvoidDeadlock(ctx, q, arg.ToAccountID, arg.FromAccountID, -arg.Amount)
		}
		return err
	})

	return result, err
}

func txnAvoidDeadlock(ctx context.Context, q *Queries, id1 int64, id2 int64, amount int64) (up1 Account, up2 Account, err error) {
	up1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     id1,
		Amount: -amount,
	})
	if err != nil {
		return
	}
	up2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     id2,
		Amount: amount,
	})
	return
}

/* But the story of deadlock is not over; the final use case is when money is being transferred
between the same two accounts in both directions. This would break our tests. So ....
1. We will create a test that will break the above solution (and I need to commit this version first!!)
2. We will fix that problem by ensuring the order of updates is in a preferred order (locks taken out
in the same order as the account ID).
To be precise:
The problem to resolve is T1 locking Acc1, T2 locking Acc2 and then each trying to lock the other account.
This will lead to deadlock. We can avoiod this problem by forcing the algorithm to lock in the order of
account ID. This means that Both T1 and T2 will attempt to lock Acc1 first (if Acc1.ID < Acc2.ID).
Only one can succeed - so deadlock cannot occur for this scenario.
First ... check this version in!!
*/

/*
So before we created a special update function AddAccountBalance ... we had this feature but it's replaced.
The AddAccountBalance avoids the need to lock-read the account to apply an update. It circumvents the issue of
a potential DEADLOCK. The deadlock occurs because we can (in the original version) have a lock on an entry
which relates to the same account record ... and so two threads are trying to lock the account. First solution
was to use to drop constraints .. biut that is ugly.
*/
// func txnAdjustAccountBalance(ctx context.Context, q *Queries, accountID int64, amount int64) (Account, error) {
// 	result, err := q.GetAccountForUpdate(ctx, accountID)
// 	if err != nil {
// 		return result, err
// 	}
// 	return q.UpdateAccount(ctx, UpdateAccountParams{
// 		ID:      accountID,
// 		Balance: result.Balance + amount,
// 	})
// }
