package models

import (
	"strings"
)

type Ips struct {
	Id    int    `db:"id"`
	Value string `db:"value"`

	Domains []Domain
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

func SliceToIpsValue(slice []string) string {
	return strings.Join(slice, ",")
}

func blocked(ip string) bool {
	return strings.HasPrefix(ip, "172") || strings.HasPrefix(ip, "192")
}

func (ips *Ips) Slice() []string {
	return strings.Split(ips.Value, ",")
}

func (ips *Ips) HasBlockedIp() bool {
	for _, ip := range ips.Slice() {
		if blocked(ip) {
			return true
		}
	}
	return false
}

func (domain *Domain) HasBlockedIps() bool {
	for _, ips := range domain.Ipss {
		if ips.HasBlockedIp() {
			return true
		}
	}
	return false
}

func (domain *Domain) UpdateIpss(ipss []Ips) (newipss []Ips) {
	appended := map[string]bool{}

	for _, ips := range domain.Ipss {
		s := ips.Value
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
		domain.Ipss = append(domain.Ipss, newipss...)
		domain.RoundRobin = len(domain.Ipss) > 1
		domain.HasBlocked = domain.HasBlockedIps()
	}
	return
}
