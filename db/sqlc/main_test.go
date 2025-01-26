package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("cannot load .env file", err)
		return
	}

	dbSource := os.Getenv("DB_URL")
	if dbSource == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	dbDriver := "postgres"

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
