package db

import (
	"database/sql"
	"log"

	"github.com/koyokr/sjb-host/models"
)

func SelectIpsWhereSlice(ips *models.Ips) (exists bool) {
	err := db.Get(
		&ips.Id,
		`SELECT id FROM ipss WHERE value=$1`,
		ips.Value,
	)
	if err == sql.ErrNoRows {
	} else if err != nil {
		log.Fatal(err)
	} else {
		exists = true
	}
	return exists
}

func SelectIps() (ipss []models.Ips) {
	err := db.Select(
		&ipss,
		`SELECT * FROM ipss`,
	)
	if err != nil {
		log.Fatal(err)
	}
	return ipss
}

func SelectIpsJoinDomainId(domainid int) (ipss []models.Ips) {
	err := db.Select(
		&ipss,
		`SELECT id, value FROM ipss AS ips
		INNER JOIN domain_to_ipss AS dti ON (dti.domain_id=$1)
		WHERE (ips.id = dti.ips_id)`,
		domainid,
	)
	if err != nil {
		log.Fatal(err)
	}
	return ipss
}

func InsertIpssWithDomainToIpss(domainid int, ipss []models.Ips) {
	tx := db.MustBegin()
	for _, ips := range ipss {
		err := tx.Get(
			&ips.Id,
			`SELECT id FROM ipss WHERE value=$1`,
			ips.Value,
		)
		if err == sql.ErrNoRows {
			err := tx.Get(
				&ips.Id,
				`INSERT INTO ipss (value) VALUES ($1) RETURNING id`,
				ips.Value,
			)
			if err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Fatal(err)
		}
		tx.MustExec(
			`INSERT INTO domain_to_ipss (domain_id,ips_id) VALUES ($1,$2)`,
			domainid, ips.Id,
		)
	}
	tx.Commit()
}
