package db

import (
	"atalkowski/go-rpc/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, accountID int64) Entry {
	entryParams := CreateEntryParams{
		AccountID: accountID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), entryParams)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.Equal(t, accountID, entry.AccountID)
	require.Equal(t, entryParams.Amount, entry.Amount)
	return entry
}

func TestCreateEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	createRandomEntry(t, account.ID)
}

func TestGetEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	entry := createRandomEntry(t, account.ID)
	entry1, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry1)
	require.Equal(t, entry.ID, entry1.ID)
	require.Equal(t, entry.AccountID, entry1.AccountID)
	require.Equal(t, entry.Amount, entry1.Amount)
	require.Equal(t, entry.CreatedAt, entry1.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account.ID)
	}
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)
	}
}
