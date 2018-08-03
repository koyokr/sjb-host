package api

import (
	"database/sql"
	"log"

	"github.com/koyokr/sjb-host/models"
)

func selectDomainWhereRoundRobinHasBlocked() (ds []models.Domain) {
	err := db.Select(
		&ds,
		`SELECT * FROM domains
		WHERE round_robin=TRUE AND has_blocked=TRUE
		ORDER BY host`,
	)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func selectDomain() (ds []models.Domain) {
	err := db.Select(
		&ds,
		`SELECT * FROM domains`,
	)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func updateDomain(d models.Domain) {
	db.MustExec(
		`UPDATE domains
		SET host=$1,round_robin=$2,has_blocked=$3
		WHERE id=$4`,
		d.Host, d.RoundRobin, d.HasBlocked, d.Id,
	)
}

func selectIpsJoinDomainId(domainid int) (ipss []models.Ips) {
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
	return
}

func insertIpssWithDomainToIpss(domainid int, ipss []models.Ips) {
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
