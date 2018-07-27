package db

import (
	"log"

	"github.com/koyokr/sjb-host/models"
)

func ExistsDomainHost(host string) (exists bool) {
	err := db.Get(
		&exists,
		`SELECT EXISTS (SELECT 1 FROM domains WHERE host=$1)`,
		host,
	)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func SelectDomainWhereRoundRobinHasBlocked() (ds []models.Domain) {
	err := db.Select(
		&ds,
		`SELECT * FROM domains
		WHERE round_robin=TRUE,has_blocked=TRUE
		ORDER BY host`,
	)
	if err != nil {
		log.Fatal(err)
	}
	return ds
}

func SelectDomain() (ds []models.Domain) {
	err := db.Select(
		&ds,
		`SELECT * FROM domains`,
	)
	if err != nil {
		log.Fatal(err)
	}
	return ds
}

func InsertDomain(d *models.Domain) {
	err := db.Get(
		&d.Id,
		`INSERT INTO domains (host,round_robin,has_blocked)
		VALUES ($1,$2,$3) RETURNING id`,
		d.Host, d.RoundRobin, d.HasBlocked,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateDomain(d models.Domain) {
	db.MustExec(
		`UPDATE domains SET host=$1,round_robin=$2,has_blocked=$3
		WHERE id=$4`,
		d.Host, d.RoundRobin, d.HasBlocked, d.Id,
	)
}
