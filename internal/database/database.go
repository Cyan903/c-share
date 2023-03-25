package database

import (
	"database/sql"
	"errors"

	"github.com/Cyan903/c-share/pkg/log"
	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB
var ErrBadPW error = errors.New("invalid password")
var ErrNotFound error = errors.New("not found")

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error.Println("Could not ping database -", err)
		return nil, err
	}

	log.Info.Println("Connected to database.")
	return db, nil
}
