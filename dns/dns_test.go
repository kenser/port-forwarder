package dns

import (
	"net"
	"testing"
)

func TestLookupIP(t *testing.T) {
	ip, err := LookupIP("no-such-host", "192.168.6.1")
	if err != nil {
		t.Error(err)
	}
	t.Log(ip)
	ips, err := net.LookupIP("www.baidu.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(ips)
}
