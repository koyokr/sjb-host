package api

import (
	"github.com/jmoiron/sqlx"
	"github.com/koyokr/sjb-host/models"
	_ "github.com/lib/pq"
)

var (
	db *sqlx.DB

	wantupdate   chan bool
	readdata     chan string
	updatedomain chan models.Domain
	updateipss   chan []models.Ips
)

func Init(conn string) (err error) {
	db, err = sqlx.Connect("postgres", conn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(20)

	wantupdate = make(chan bool)
	readdata = make(chan string)
	updatedomain = make(chan models.Domain)
	updateipss = make(chan []models.Ips)

	go controlDataLoop()
	go updateDataLoop()

	return nil
}

func Close() {
	close(wantupdate)
	close(readdata)
	close(updatedomain)
	close(updateipss)
	db.Close()
}
