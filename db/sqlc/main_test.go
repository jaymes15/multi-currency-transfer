package db

import (
	"database/sql"
	"lemfi/simplebank/config"
	"lemfi/simplebank/db"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// Load environment variables
	godotenv.Load()

	// Set configuration
	config.Set()

	// Connect to database using the same approach as bootstrap
	db.Connect()
	testDB = db.GetPostgresDBConnection()

	// Create queries instance
	testQueries = New(testDB)

	os.Exit(m.Run())
}
