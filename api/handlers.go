package api

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/koyokr/sjb-host/db"
)

func GetDomainAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var s strings.Builder
	for _, d := range db.SelectDomainWhereRoundRobinHasBlocked() {
		s.WriteString(d.Host)
		s.WriteString(":")
		for _, ips := range db.SelectIpsJoinDomainId(d.Id) {
			s.WriteString("\n\t")
			s.WriteString(ips.Value)
		}
		s.WriteString("\n\n")
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, s.String())
}

func PutDomainFunc(ch chan string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := ps.ByName("name")
		if net.ParseIP(name) != nil {
			NotFound(w, r)
			return
		}
		if name == "favicon.ico" || name == "robots.txt" {
			NotFound(w, r)
			return
		}

		exists := db.ExistsDomainHost(name)
		if exists {
			fmt.Fprintf(w, "Exists")
			return
		}
		ch <- name

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Put")
	}
}
