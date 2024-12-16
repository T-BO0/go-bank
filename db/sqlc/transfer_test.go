package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/T-BO0/bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	args := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, args.FromAccountID)
	require.Equal(t, transfer.ToAccountID, args.ToAccountID)
	require.Equal(t, transfer.Amount, args.Amount)

	return transfer
}

func createRandomTransferForAccounts(t *testing.T, fromAccountId int64, toAccountId int64) Transfer {
	args := CreateTransferParams{
		FromAccountID: fromAccountId,
		ToAccountID:   toAccountId,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, args.FromAccountID)
	require.Equal(t, transfer.ToAccountID, args.ToAccountID)
	require.Equal(t, transfer.Amount, args.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestGetTransferByAccounts(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)
	transfer1 := createRandomTransferForAccounts(t, fromAcc.ID, toAcc.ID)

	args := GetTransferByAccountsParams{FromAccountID: transfer1.FromAccountID, ToAccountID: transfer1.ToAccountID}
	transfer2, err := testQueries.GetTransferByAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer2.ID, transfer1.ID)
	require.Equal(t, transfer2.Amount, transfer1.Amount)
	require.Equal(t, transfer2.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer2.ToAccountID, transfer1.ToAccountID)
	require.WithinDuration(t, transfer2.CreatedAt, transfer1.CreatedAt, time.Second)
}

func TestGetTransferByFromAccount(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)

	transfer1 := createRandomTransferForAccounts(t, fromAcc.ID, toAcc.ID)

	transfer2, err := testQueries.GetTransferByFromAccountId(context.Background(), transfer1.FromAccountID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer2.ID, transfer1.ID)
	require.Equal(t, transfer2.Amount, transfer1.Amount)
	require.Equal(t, transfer2.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer2.ToAccountID, transfer1.ToAccountID)
	require.WithinDuration(t, transfer2.CreatedAt, transfer1.CreatedAt, time.Second)
}

func TestGetTransferByToAccount(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)

	transfer1 := createRandomTransferForAccounts(t, fromAcc.ID, toAcc.ID)

	transfer2, err := testQueries.GetTransferByToAccountId(context.Background(), transfer1.ToAccountID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer2.ID, transfer1.ID)
	require.Equal(t, transfer2.Amount, transfer1.Amount)
	require.Equal(t, transfer2.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer2.ToAccountID, transfer1.ToAccountID)
	require.WithinDuration(t, transfer2.CreatedAt, transfer1.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestListTransfer(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	args := ListTransferParams{Limit: 5, Offset: 5}
	transfers, err := testQueries.ListTransfer(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestListTransferByAccounts(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransferForAccounts(t, fromAcc.ID, toAcc.ID)
	}

	args := ListTransferByAccountsParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransferByAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, fromAcc.ID)
		require.Equal(t, transfer.ToAccountID, toAcc.ID)
	}
}

func TestListTransferByFromAccountId(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransferForAccounts(t, fromAcc.ID, toAcc.ID)
	}

	args := ListTransferByFromAccountIdParams{
		FromAccountID: fromAcc.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransferByFromAccountId(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, fromAcc.ID)
	}
}

func TestListTransferByToAccountId(t *testing.T) {
	fromAcc := createRandomAccount(t)
	toAcc := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransferForAccounts(t, fromAcc.ID, toAcc.ID)
	}

	args := ListTransferByToAccountIdParams{
		ToAccountID: toAcc.ID,
		Limit:       5,
		Offset:      5,
	}

	transfers, err := testQueries.ListTransferByToAccountId(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.ToAccountID, toAcc.ID)
	}
}

func TestUpdateTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	args := UpdateTransferParams{ID: transfer1.ID, Amount: util.RandomMoney()}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer2.ID, transfer1.ID)
	require.Equal(t, transfer2.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer2.ToAccountID, transfer1.ToAccountID)
	require.Equal(t, transfer2.Amount, args.Amount)
	require.WithinDuration(t, transfer2.CreatedAt, transfer1.CreatedAt, time.Second)
}
