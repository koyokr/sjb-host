package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/koyokr/sjb-host/models"
)

func GetDomains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, readData())
}

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

func readData() string {
	wantupdate <- false
	return <-readdata
}

func updateData(domain models.Domain, ipss []models.Ips) {
	wantupdate <- true
	updatedomain <- domain
	updateipss <- ipss
}

func controlDataLoop() {
	data := createResult()
	for {
		if <-wantupdate {
			domain := <-updatedomain
			ipss := <-updateipss

			updateDomain(domain)
			insertIpssWithDomainToIpss(domain.Id, ipss)
			data = createResult()
		} else {
			readdata <- data
		}
	}
}

func updateDataLoop() {
	for {
		for _, d := range selectDomain() {
			d.Ipss = selectIpsJoinDomainId(d.Id)

			ipss := lookupAll(d.Host)
			newipss := d.UpdateIpss(ipss)

			if len(newipss) != 0 {
				updateData(d, newipss)
			}
			time.Sleep(5 * time.Second)
		}
		time.Sleep(1 * time.Hour)
	}
}
