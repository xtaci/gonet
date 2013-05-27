package geoip

import (
	"bufio"
	"net"
	"os"
	"strings"
)

import (
	itree "misc/alg/interval_tree"
)

const (
	GEOIP = "GeoIPCountryWhois.csv"
)

var _geoip itree.Tree

func init() {
	path := os.Getenv("GOPATH") + "/src/misc/geoip/" + GEOIP
	file, err := os.Open(path)
	if err != nil {
		panic("error opening geoip file")
	}

	r := bufio.NewReader(file)
	for {
		line, e := r.ReadString('\n')
		line = strings.TrimSpace(line)

		// empty-line & #comment
		if line == "" {
			if e == nil {
				continue
			} else {
				break
			}
		}

		// split fields
		fields := strings.Split(line, ",")
		from := net.ParseIP(strings.Trim(fields[0], `"`))
		to := net.ParseIP(strings.Trim(fields[1], `"`))

		ifrom := _int64_ip(from)
		ito := _int64_ip(to)

		for i := 0; i < len(fields); i++ {
			_geoip.Insert(ifrom, ito, strings.Trim(fields[4], `"`))
		}
	}
}

func _int64_ip(ip net.IP) int64 {
	idx := 0
	if len(ip) > 4 {
		idx = 12
	}

	ipv4 := uint32(ip[0+idx])<<24 | uint32(ip[1+idx])<<16 | uint32(ip[2+idx])<<8 | uint32(ip[3+idx])
	return int64(ipv4)
}

//------------------------------------------------ Get Country Code
func Query(ip net.IP) string {
	i64 := _int64_ip(ip)
	if n := _geoip.Lookup(i64, i64); n != nil {
		return n.Data().(string)
	}

	return ""
}
