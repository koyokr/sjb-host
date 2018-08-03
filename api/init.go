package api

import (
	"github.com/jmoiron/sqlx"
	"github.com/koyokr/sjb-host/models"
	_ "github.com/lib/pq"
)

var (
	db *sqlx.DB

	wantupdate   = make(chan bool)
	readdata     = make(chan string)
	updatedomain = make(chan models.Domain)
	updateipss   = make(chan []models.Ips)
)

func InitDB(conn string) (err error) {
	db, err = sqlx.Connect("postgres", conn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(20)

	go controlDataLoop()
	go updateDataLoop()
	return
}

func CloseDB() {
	db.Close()
}
