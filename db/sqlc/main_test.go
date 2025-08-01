package db

import (
	"lemfi/simplebank/config"
	"lemfi/simplebank/db"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var testQueries Store
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	// Load environment variables
	godotenv.Load()

	// Set configuration
	config.Set()

	// Connect to database using the same approach as bootstrap
	db.Connect()
	testDB = db.GetPostgresDBConnection()

	// Create queries instance
	testQueries = NewStore(testDB)

	os.Exit(m.Run())
}
