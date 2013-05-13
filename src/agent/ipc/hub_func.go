package ipc

import (
	"time"
	"log"
)

import (
	"hub/protos"
	"misc/packet"
)

type Info struct {
	Id int32
	State int32
	Score int32
	Clan int32
	ProtectTime int64
	Name string
}

func Login(id int32) bool {
	defer _hub_err()
	req := protos.ID{}
	req.F_id = id
	ret := _call(packet.Pack(protos.Code["login_req"], req, nil))
	reader := packet.Reader(ret)
	tbl, err := protos.PKT_INT(reader)

	if err != nil || tbl.F_v==0 {
		return false
	}

	return true
}

func Logout(id int32) bool {
	defer _hub_err()
	req := protos.ID{}
	req.F_id = id
	ret := _call(packet.Pack(protos.Code["logout_req"], req, nil))
	reader := packet.Reader(ret)
	tbl, err := protos.PKT_INT(reader)

	if err != nil || tbl.F_v==0 {
		return false
	}

	return true
}

func Raid(id int32) bool {
	defer _hub_err()

	req := protos.ID{}
	req.F_id = id
	ret := _call(packet.Pack(protos.Code["raid_req"], req, nil))
	reader := packet.Reader(ret)
	tbl, err := protos.PKT_INT(reader)

	if err != nil || tbl.F_v==0 {
		return false
	}

	return true
}

func Protect(id int32, until time.Time) bool {
	defer _hub_err()

	req := protos.ID{}
	req.F_id = id
	ret := _call(packet.Pack(protos.Code["protect_req"], req, nil))
	reader := packet.Reader(ret)
	tbl, err := protos.PKT_INT(reader)

	if err != nil || tbl.F_v==0 {
		return false
	}

	return true
}

func Free(id int32) bool {
	defer _hub_err()

	req := protos.ID{}
	req.F_id = id
	ret := _call(packet.Pack(protos.Code["free_req"], req, nil))
	reader := packet.Reader(ret)
	tbl, err := protos.PKT_INT(reader)

	if err != nil || tbl.F_v==0 {
		return false
	}

	return true
}

func Unprotect(id int32) bool {
	defer _hub_err()

	req := protos.ID{}
	req.F_id = id
	ret := _call(packet.Pack(protos.Code["unprotect_req"], req, nil))
	reader := packet.Reader(ret)
	tbl, err := protos.PKT_INT(reader)

	if err != nil || tbl.F_v==0 {
		return false
	}

	return true
}

func GetInfo(id int32) (info Info,err error) {
	defer _hub_err()

	req := protos.ID{}
	req.F_id = id
	ret := _call(packet.Pack(protos.Code["getinfo_req"], req, nil))
	reader := packet.Reader(ret)
	tbl, _err := protos.PKT_INFO(reader)
	info.Id = tbl.F_id
	info.State = tbl.F_state
	info.Score = tbl.F_score
	info.Clan = tbl.F_clan
	info.ProtectTime = tbl.F_protecttime
	info.Name = tbl.F_name

	return info, _err
}

func GetList(A,B int32) (ids, scores []int32, err error) {
	defer _hub_err()

	req := protos.GETLIST{}
	req.F_A = A
	req.F_B = B
	ret := _call(packet.Pack(protos.Code["getlist_req"], req, nil))
	reader := packet.Reader(ret)
	tbl, _err := protos.PKT_LIST(reader)
	ids = make([]int32, len(ids))
	scores = make([]int32, len(scores))

	for k := range tbl.F_items {
		ids[k] = tbl.F_items[k].F_id
		scores[k] = tbl.F_items[k].F_score
	}

	return ids, scores, _err
}

func _hub_err() {
	if x := recover(); x != nil {
		log.Println(x)
	}
}
