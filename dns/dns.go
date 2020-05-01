package dns

import (
	"errors"
	"github.com/miekg/dns"
)

func getIp(domain string) string {
	//net.LookupIP()
	return ""
}

var (
	errNoSuchHost = errors.New("no such host")
)

// DNSError represents a DNS lookup error.
type DNSError struct {
	Err         string // description of the error
	Name        string // name looked for
	Server      string // server used
	IsTimeout   bool   // if true, timed out; not all timeouts set this
	IsTemporary bool   // if true, error is temporary; not all errors set this
	IsNotFound  bool   // if true, host could not be found
}

func (e *DNSError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := "lookup " + e.Name
	if e.Server != "" {
		s += " on " + e.Server
	}
	s += ": " + e.Err
	return s
}

func LookupIP(host string, resolver string) (ip string, err error) {
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(host+".", dns.TypeA)
	r, _, err := c.Exchange(&m, resolver+":53")
	if err != nil {
		return ip, err
	}
	if len(r.Answer) == 0 {
		return ip, &DNSError{
			Err:    errNoSuchHost.Error(),
			Name:   host,
			Server: resolver,
		}
	}
	for _, ans := range r.Answer {
		switch t := ans.(type) {
		case *dns.A:
			return t.A.String(), nil
		default:
			continue
		}
	}
	return ip, &DNSError{
		Err:    errNoSuchHost.Error(),
		Name:   host,
		Server: resolver,
	}
}
