package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func New() *sql.DB {
	connStr := "host=database port=5432 user=postgres dbname=postgres sslmode=disable password=password"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}
