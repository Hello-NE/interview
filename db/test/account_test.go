package db

import (
	"context"
	db "interview/db/sqlc"
	util "interview/db/util"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestTransferTx(t *testing.T) {
	testStore := db.NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan db.TransferTxResult)
	log.Printf("account1 balance before: %d, account2 balance before: %d", account1.Balance, account2.Balance)
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), db.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)
		require.Equal(t, account1.ID, result.FromAccount.ID)
		require.Equal(t, account2.ID, result.ToAccount.ID)
		require.Equal(t, amount, result.FromEntry.Amount)
		require.Equal(t, -amount, result.ToEntry.Amount)
		log.Printf("transfer %d: from account %d to account %d, amount: %d", i+1, result.FromAccount.ID, result.ToAccount.ID, amount)
		log.Printf("account1 balance after: %d, account2 balance after: %d", result.FromAccount.Balance, result.ToAccount.Balance)
	}

}

func createRandomAccount(t *testing.T) db.Account {
	args := db.CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), args)
	if err != nil {
		t.Fatal("failed to create account:", err)
	}
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
