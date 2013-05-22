package cd_client

import (
	"time"
	"log"
)

import (
	cd "cooldown/protos"
	"misc/packet"
)

func Add(oid int32, user_id int32, timeout int64) bool {
	defer _cd_err()
	req := cd.ADD_REQ{}
	req.F_oid = oid
	req.F_user_id = user_id
	req.F_timeout = timeout
	ret := _call(packet.Pack(cd.Code["add_req"], req, nil))
	reader := packet.Reader(ret)
	tbl, err := hub.PKT_INT(reader)

	if err != nil || tbl.F_v==0 {
		return false
	}

	return true
}

func _cd_err() {
	if x := recover(); x != nil {
		log.Println(x)
	}
}
