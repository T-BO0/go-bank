package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/T-BO0/bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)

	args := CreateEntryParams{AccountID: account.ID, Amount: util.RandomMoney()}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, account.ID)
	require.Equal(t, entry.Amount, args.Amount)

	return entry
}

func createRandomEntryForAccount(t *testing.T, account Account) Entry {

	args := CreateEntryParams{AccountID: account.ID, Amount: util.RandomMoney()}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, account.ID)
	require.Equal(t, entry.Amount, args.Amount)

	return entry
}

func compareTwoEntries(t *testing.T, err error, entry2 Entry, entry1 Entry) {
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry2.ID, entry1.ID)
	require.Equal(t, entry2.AccountID, entry1.AccountID)
	require.Equal(t, entry2.Amount, entry1.Amount)
	require.WithinDuration(t, entry2.CreatedAt, entry1.CreatedAt, time.Second)
}

func TestCreateEntrie(t *testing.T) {
	createRandomEntry(t)
}

func TestCreateEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	compareTwoEntries(t, err, entry2, entry1)
}

func TestGetEntryByAccountId(t *testing.T) {
	entry1 := createRandomEntry(t)

	entry2, err := testQueries.GetEntryByAccountId(context.Background(), entry1.AccountID)
	compareTwoEntries(t, err, entry2, entry1)
}

func TestDeleteEntrie(t *testing.T) {
	entry1 := createRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListEntry(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	args := ListEntryParams{Limit: 5, Offset: 5}
	entries, err := testQueries.ListEntry(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestListEntryByAccountId(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntryForAccount(t, account)
	}

	args := ListEntryByAccountIdParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entryes, err := testQueries.ListEntryByAccountId(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entryes, 5)

	for _, entry := range entryes {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, account.ID)
	}
}

func TestUpdateEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	args := UpdateEntryParams{ID: entry1.ID, Amount: util.RandomMoney()}

	entry2, err := testQueries.UpdateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry2.Amount, args.Amount)
}
