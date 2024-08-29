package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

func CheckMailHash(email, hash string) (bool, error) {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return false, err
	}
	defer db.Close()

	cmd := "SELECT hash FROM users WHERE email=?"
	row := db.QueryRow(cmd, email)
	var h string
	err = row.Scan(&h)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if hash != h {
		return false, nil
	}

	return true, nil
}

func GetDisplayName(email string) (string, error) {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return "", err
	}
	defer db.Close()

	cmd := "SELECT displayName FROM users WHERE email=?"
	row := db.QueryRow(cmd, email)
	var displayName string
	err = row.Scan(&displayName)
	if err != nil {
		return "", err
	}

	return displayName, nil
}

func CheckExistMail(email string) (bool, error) {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return false, err
	}
	defer db.Close()

	cmd := "SELECT 1 FROM users WHERE email=?"
	row := db.QueryRow(cmd, email)
	var exists bool
	err = row.Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func CheckExistVerify(email string) (bool, error) {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return false, err
	}
	defer db.Close()

	cmd := "SELECT 1 FROM verify WHERE email=?"
	row := db.QueryRow(cmd, email)
	var exists bool
	err = row.Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func UserInsert(email, hash, displayName string, isOauth bool) error {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	cmd := "INSERT INTO users (email,oauth,hash,displayName) VALUES (?,?,?,?)"
	_, err = db.Exec(cmd, email, isOauth, hash, displayName)
	if err != nil {
		return err
	}

	return nil
}

func VerifyInsert(email, hash, displayName, verifyCode string) error {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	cmd := "INSERT INTO verify (email,hash,displayName,code) VALUES (?,?,?,?)"
	_, err = db.Exec(cmd, email, hash, displayName, verifyCode)
	if err != nil {
		return err
	}

	return nil
}

func VerifyCodeCmp(email, code string) (bool, error) {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return false, err
	}
	defer db.Close()

	cmd := "SELECT code FROM verify WHERE email=?"
	row := db.QueryRow(cmd, email)
	var c string
	err = row.Scan(&c)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if c != code {
		return false, nil
	}

	return true, nil
}

func MoveData(email string) error {
	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	cmd := "SELECT hash,displayName FROM verify WHERE email=?"
	row := db.QueryRow(cmd, email)
	var hash string
	var displayName string
	err = row.Scan(&hash, &displayName)
	if err != nil {
		return err
	}
	cmd = "DELETE FROM verify WHERE email=?"
	_, err = db.Exec(cmd, email)
	if err != nil {
		return err
	}
	UserInsert(email, hash, displayName, false)

	return nil
}
