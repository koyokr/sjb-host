package api

import (
	"strings"
	"time"

	"github.com/koyokr/sjb-host/models"
)

func createResult() string {
	var s strings.Builder
	for _, domain := range selectDomainWhereRoundRobinHasBlocked() {
		s.WriteString(domain.Host)
		for _, ips := range selectIpsJoinDomainId(domain.Id) {
			s.WriteString("\n\t")
			s.WriteString(ips.Value)
		}
		s.WriteString("\n\n")
	}
	return s.String()
}

func updateAsyncDB() (chan bool, chan string) {
	is_update := make(chan bool)
	read_data := make(chan string)
	write_domain := make(chan models.Domain)
	write_ipss := make(chan []models.Ips)

	go func() {
		data := createResult()
		for {
			if <-is_update {
				domain := <-write_domain
				ipss := <-write_ipss

				updateDomain(domain)
				insertIpssWithDomainToIpss(domain.Id, ipss)
				data = createResult()
			} else {
				read_data <- data
			}
		}
	}()

	go func() {
		for {
			for _, d := range selectDomain() {
				time.Sleep(4 * time.Second)
				d.Ipss = selectIpsJoinDomainId(d.Id)

				ipss := lookupAll(d.Host)
				newipss := d.UpdateIpss(ipss)
				if len(newipss) == 0 {
					continue
				}

				is_update <- true
				write_domain <- d
				write_ipss <- newipss
			}
			time.Sleep(1 * time.Hour)
		}
	}()

	return is_update, read_data
}
