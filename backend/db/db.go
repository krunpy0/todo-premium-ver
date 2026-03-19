package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(connString string) error {
	var err error
	DB, err = sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	return DB.Ping()
}