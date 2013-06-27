package main

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"time"
)

import (
	"cfg"
	"helper"
)

const (
	DEFAULT_BAN_TIME = 5
)

//---------------------------------------------------------- IP->UnBan time
var _banned_ips map[uint32]int64

func init() {
	_banned_ips = make(map[uint32]int64)
}

//---------------------------------------------------------- ban an ip
func Ban(_ip net.IP) {
	config := cfg.Get()
	ban_time, err := strconv.Atoi(config["ban_time"])
	if err != nil {
		ban_time = DEFAULT_BAN_TIME
		log.Println("cannot get ban_timeout from config", err)
	}

	intip := _ip2int(_ip)

	// randomize the timeout, for effective DoS protection
	ban := uint32(ban_time)
	_banned_ips[intip] = time.Now().Unix() + int64(ban+helper.LCG()%ban)
}

//---------------------------------------------------------- test whether the ip is banned
func IsBanned(_ip net.IP) bool {
	intip := _ip2int(_ip)
	timeout, exists := _banned_ips[intip]

	if !exists {
		return false
	} else if timeout < time.Now().Unix() {
		delete(_banned_ips, intip)
		return false
	} else {
		return true
	}
}

func _ip2int(_ip net.IP) uint32 {
	ip := _ip.To4()
	if ip != nil {
		return binary.BigEndian.Uint32(ip)
	}

	return 0
}
