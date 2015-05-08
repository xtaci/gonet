package geoip

import (
	"net"
	"testing"
)

func TestGeoIP(t *testing.T) {
	ip := net.ParseIP("103.14.100.100")
	t.Log(Query(ip))
	ip = net.ParseIP("171.216.90.18")
	t.Log(Query(ip))
	ip = net.ParseIP("8.35.201.33")
	t.Log(Query(ip))
	ip = net.ParseIP("62.216.125.241")
	t.Log(Query(ip))
}

func BenchmarkGeoIP(b *testing.B) {
	ip := net.ParseIP("103.14.100.100")

	for i := 0; i < b.N; i++ {
		Query(ip)
	}
}
