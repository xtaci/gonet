package geoip

import (
	"bufio"
	"encoding/binary"
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

		if line == "" {
			if e == nil {
				continue
			} else {
				break
			}
		}

		fields := strings.Split(line, ",")
		from := net.ParseIP(strings.Trim(fields[0], `"`))
		to := net.ParseIP(strings.Trim(fields[1], `"`))

		ifrom := _int64_ip(from)
		ito := _int64_ip(to)

		if ifrom*ito != 0 {
			_geoip.Insert(ifrom, ito, strings.Trim(fields[4], `"`))
		}
	}
}

func _int64_ip(_ip net.IP) int64 {
	ip := _ip.To4()
	if ip != nil {
		ipv4 := binary.BigEndian.Uint32(ip)
		return int64(ipv4)
	}

	return 0
}

//------------------------------------------------ Get Country Code
func Query(ip net.IP) string {
	i64 := _int64_ip(ip)
	if n := _geoip.Lookup(i64, i64); n != nil {
		return n.Data().(string)
	}

	return ""
}
