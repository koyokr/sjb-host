package models

import (
	"strings"
)

type Ips struct {
	Id    int    `db:"id"`
	Value string `db:"value"`

	Domains []Domain
}

func SliceToIpsValue(slice []string) string {
	return strings.Join(slice, ",")
}

func (ips *Ips) Slice() []string {
	return strings.Split(ips.Value, ",")
}

type DomainToIps struct {
	DomainId int `db:"domain_id"`
	IpsId    int `db:"ips_id"`
}

type Domain struct {
	Id         int    `db:"id"`
	Host       string `db:"host"`
	RoundRobin bool   `db:"round_robin"`
	HasBlocked bool   `db:"has_blocked"`

	Ipss []Ips
}

func blocked(ip string) bool {
	return strings.HasPrefix(ip, "172") || strings.HasPrefix(ip, "192")
}

func containsBlocked(ipss []Ips) bool {
	for _, ips := range ipss {
		for _, ip := range ips.Slice() {
			if blocked(ip) {
				return true
			}
		}
	}
	return false
}

func (d *Domain) UpdateIpss(ipss []Ips) (newipss []Ips) {
	appended := map[string]bool{}

	for _, dips := range d.Ipss {
		s := dips.Value
		appended[s] = true
	}

	for _, ips := range ipss {
		s := ips.Value
		if appended[s] != true {
			appended[s] = true
			newipss = append(newipss, ips)
		}
	}

	if len(newipss) > 0 {
		d.Ipss = append(d.Ipss, newipss...)
		d.RoundRobin = len(d.Ipss) > 1
		d.HasBlocked = containsBlocked(d.Ipss)
	}
	return newipss
}
