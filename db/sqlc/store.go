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

// Transfer performs all necessary updates to execute a money transfer between 2 accounts.
// It creates a transfer record, adds account entries and updates the accont balances for each account
// Within a single transcation
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// Note use of in line function here:-
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: update accounts' balances ... but concern over locking... so we will come back to this
		// result.FromAccount, err = q.UpdateAccount(ctx, )
		return nil
	})

	return result, err
}
