package db

import "database/sql"

func GetPostgresDBConnection() *sql.DB {
	return PostgresDB
}
