package db

import (
	"database/sql"
	db "interview/db/sqlc"
	"log"
	"os"
	"testing"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *db.Queries
var testDB *sql.DB
func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testDB = conn
	testQueries = db.New(conn)
	os.Exit(m.Run())
}
