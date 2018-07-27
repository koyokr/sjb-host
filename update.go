package main

import (
	"time"

	"github.com/koyokr/sjb-host/db"
	"github.com/koyokr/sjb-host/models"
)

func updateDomainAll() {
	for _, d := range db.SelectDomain() {
		time.Sleep(4 * time.Second)
		d.Ipss = db.SelectIpsJoinDomainId(d.Id)

		ipss := lookupAll(d.Host)
		newipss := d.UpdateIpss(ipss)
		if len(newipss) == 0 {
			continue
		}
		db.UpdateDomain(d)
		db.InsertIpssWithDomainToIpss(d.Id, newipss)
	}
	time.Sleep(1 * time.Hour)
}

func insertDomainFromChan(ch chan string) {
	d := models.Domain{Host: <-ch}

	ipss := lookupAll(d.Host)
	newipss := d.UpdateIpss(ipss)
	if len(newipss) == 0 {
		return
	}
	db.InsertDomain(&d)
	db.InsertIpssWithDomainToIpss(d.Id, newipss)
}
