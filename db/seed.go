package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Seed() (bool, error) {
	db, err := sql.Open("sqlite3", "../database.sqlite")
	if err != nil {
		return false, err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return false, err
	}
	return false, nil
}
