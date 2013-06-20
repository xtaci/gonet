package geoip

import (
	"net"
	"testing"
)

func TestGeoIP(t *testing.T) {
	ip := net.ParseIP("103.14.100.100")

	t.Log(Query(ip))
	if Query(ip) != "HK" {
		t.Error("ip mismatch")
	}
}

func BenchmarkGeoIP(b *testing.B) {
	ip := net.ParseIP("103.14.100.100")

	for i := 0; i < b.N; i++ {
		Query(ip)
	}
}
