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

func updateAsyncDB() (is_update chan bool, read_data chan string) {
	var (
		write_domain chan models.Domain
		write_ipss   chan []models.Ips
	)

	go func() {
		var data string
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
