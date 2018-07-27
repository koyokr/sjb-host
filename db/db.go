package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db *sqlx.DB
)

func InitDB(conn string) (err error) {
	db, err = sqlx.Connect("postgres", conn)
	db.SetMaxOpenConns(20)
	return err
}

func CloseDB() {
	db.Close()
}
