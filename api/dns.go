package api

import (
	"sort"

	"github.com/koyokr/sjb-host/models"
	"github.com/miekg/dns"
)

func lookup(cli *dns.Client, host string, ns string) (ips models.Ips) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(host), dns.TypeA)
	r, _, err := cli.Exchange(msg, ns)
	if err != nil && r.Rcode != dns.RcodeSuccess {
		return
	}

	var slice []string
	for _, a := range r.Answer {
		if t, ok := a.(*dns.A); ok {
			slice = append(slice, t.A.String())
		}
	}
	sort.Strings(slice)
	ips.Value = models.SliceToIpsValue(slice)
	return
}

func lookupAll(host string) (ipss []models.Ips) {
	nss := [...]string{
		"1.1.1.1:53",        // Cloudflare
		"8.8.8.8:53",        // Google
		"9.9.9.9:53",        // Quad9
		"208.67.222.222:53", // OpenDNS
	}

	cli := new(dns.Client)
	ch := make(chan models.Ips)
	defer close(ch)

	for _, ns := range nss {
		go func(ns string) {
			ch <- lookup(cli, host, ns)
		}(ns)
	}
	for range nss {
		ips := <-ch
		if ips.Value != "" {
			ipss = append(ipss, ips)
		}
	}
	return
}
