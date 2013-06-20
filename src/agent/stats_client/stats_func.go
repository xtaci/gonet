package stats_client

import (
	"log"
)

import (
	"misc/packet"
	stats "stats/protos"
)

func Ping() bool {
	defer _stats_err()
	req := stats.INT{}
	req.F_v = 1
	ret := _call(packet.Pack(stats.Code["ping_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, _ := stats.PKT_INT(reader)
	if tbl.F_v != req.F_v {
		return false
	}

	return true
}

func _stats_err() {
	if x := recover(); x != nil {
		log.Println(x)
	}
}
