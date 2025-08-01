package db

import "github.com/jackc/pgx/v5/pgxpool"

func GetPostgresDBConnection() *pgxpool.Pool {
	return PostgresDB
}
