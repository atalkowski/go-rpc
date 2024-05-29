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

		// Adjust the FromAccount balance (note the use of SELECT ... FOR UPDATE in this call)
		result.FromAccount, err = txnAdjustAccountBalance(ctx, q, arg.FromAccountID, -arg.Amount)
		if err != nil {
			return err
		}

		// Adjust the ToAccount balance (note the use of SELECT ... FOR UPDATE in this call)
		result.ToAccount, err = txnAdjustAccountBalance(ctx, q, arg.ToAccountID, arg.Amount)
		return err
	})

	return result, err
}

func txnAdjustAccountBalance(ctx context.Context, q *Queries, accountID int64, amount int64) (Account, error) {
	result, err := q.GetAccountForUpdate(ctx, accountID)
	if err != nil {
		return result, err
	}
	return q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      accountID,
		Balance: result.Balance + amount,
	})
}
