package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	db "interview/db/sqlc"
	"interview/util"
)

var testQueries *db.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// runDBMigration("../migration", config.DBDriver, config.DBSource)

	testQueries = db.New(testDB)

	os.Exit(m.Run())
}


func createRandomUser(t *testing.T) db.User {
	hashedPassword := util.RandomString(6)
	user, err := testQueries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     util.RandomOwner(),
		PasswordHash: hashedPassword,
		FullName:     util.RandomOwner(),
		Email:        util.RandomEmail(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.PasswordHash, hashedPassword)

	return user
}
