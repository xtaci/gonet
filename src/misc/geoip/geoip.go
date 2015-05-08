package geoip

import (
	"encoding/binary"
	"encoding/csv"
	"log"
	"net"
	"os"
	"strconv"
)

import (
	itree "misc/alg/interval_tree"
)

const (
	GEOIP            = "GeoIPCountryWhois.csv"
	FROM_IDX         = 2
	TO_IDX           = 3
	COUNTRY_CODE_IDX = 4
)

var _geoip itree.Tree

func init() {
	log.Println("Loading GEOIP...")
	defer log.Println("GEOIP Load Complete.")
	path := os.Getenv("GOPATH") + "/src/misc/geoip/" + GEOIP
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		panic("error opening geoip file")
	}
	defer file.Close()

	csv_reader := csv.NewReader(file)
	records, err := csv_reader.ReadAll()
	if err != nil {
		log.Println(err)
		panic("cannot read file geoip file")
	}

	for k := range records {
		from, _ := strconv.Atoi(records[k][FROM_IDX])
		to, _ := strconv.Atoi(records[k][TO_IDX])

		if from <= to {
			_geoip.Insert(int64(from), int64(to), records[k][COUNTRY_CODE_IDX])
		}
	}
}

//------------------------------------------------ Get Country Code
func Query(ip net.IP) string {
	i64 := _int64_ip(ip)
	if n := _geoip.Lookup(i64, i64); n != nil {
		return n.Data().(string)
	}

	return ""
}

func _int64_ip(_ip net.IP) int64 {
	ip := _ip.To4()
	if ip != nil {
		ipv4 := binary.BigEndian.Uint32(ip)
		return int64(ipv4)
	}

	return 0
}
