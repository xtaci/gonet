package hub_client

import (
	"log"
	"time"
)

import (
	hub "hub/protos"
	"misc/packet"
	. "types"
)

type Info struct {
	Id          int32
	State       byte
	Score       int32
	ProtectTime int64
}

func Ping() bool {
	defer _hub_err()
	req := hub.INT{}
	req.F_v = 1
	ret := _call(packet.Pack(hub.Code["ping_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, _ := hub.PKT_INT(reader)
	if tbl.F_v != req.F_v {
		return false
	}

	return true
}

func Login(id int32) bool {
	defer _hub_err()
	req := hub.LOGIN_REQ{}
	req.F_id = id
	ret := _call(packet.Pack(hub.Code["login_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, err := hub.PKT_LOGIN_ACK(reader)

	if err != nil || !tbl.F_success {
		return false
	}

	return true
}

func Logout(id int32) bool {
	defer _hub_err()
	req := hub.ID{}
	req.F_id = id
	ret := _call(packet.Pack(hub.Code["logout_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, err := hub.PKT_INT(reader)

	if err != nil || tbl.F_v == 0 {
		return false
	}

	return true
}

func Raid(id int32) bool {
	defer _hub_err()

	req := hub.ID{}
	req.F_id = id
	ret := _call(packet.Pack(hub.Code["raid_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, err := hub.PKT_INT(reader)

	if err != nil || tbl.F_v == 0 {
		return false
	}

	return true
}

func Protect(id int32, until time.Time) bool {
	defer _hub_err()

	req := hub.ID{}
	req.F_id = id
	ret := _call(packet.Pack(hub.Code["protect_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, err := hub.PKT_INT(reader)

	if err != nil || tbl.F_v == 0 {
		return false
	}

	return true
}

func Free(id int32) bool {
	defer _hub_err()

	req := hub.ID{}
	req.F_id = id
	ret := _call(packet.Pack(hub.Code["free_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, err := hub.PKT_INT(reader)

	if err != nil || tbl.F_v == 0 {
		return false
	}

	return true
}

func Unprotect(id int32) bool {
	defer _hub_err()

	req := hub.ID{}
	req.F_id = id
	ret := _call(packet.Pack(hub.Code["unprotect_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, err := hub.PKT_INT(reader)

	if err != nil || tbl.F_v == 0 {
		return false
	}

	return true
}

func GetInfo(id int32) (info Info, flag bool) {
	defer _hub_err()

	req := hub.ID{}
	req.F_id = id
	ret := _call(packet.Pack(hub.Code["getinfo_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, _err := hub.PKT_INFO(reader)
	if _err == nil && tbl.F_flag {
		info.Id = tbl.F_id
		info.State = tbl.F_state
		info.Score = tbl.F_score
		info.ProtectTime = tbl.F_protecttime
		flag = true
		return
	}

	flag = false
	return
}

func GetList(A, B int32) (ids, scores []int32, err error) {
	defer _hub_err()

	req := hub.GETLIST{}
	req.F_A = A
	req.F_B = B
	ret := _call(packet.Pack(hub.Code["getlist_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, _err := hub.PKT_LIST(reader)
	ids = make([]int32, len(ids))
	scores = make([]int32, len(scores))

	for k := range tbl.F_items {
		ids[k] = tbl.F_items[k].F_id
		scores[k] = tbl.F_items[k].F_score
	}

	return ids, scores, _err
}

func AddUser(id int32) bool {
	defer _hub_err()
	req := hub.ID{}
	req.F_id = id
	ret := _call(packet.Pack(hub.Code["adduser_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, err := hub.PKT_INT(reader)

	if err != nil || tbl.F_v == 0 {
		return false
	}

	return true
}

//------------------------------------------------ Forward to Hub
func Forward(req *IPCObject) bool {
	defer _hub_err()
	// HUB protocol forwarding
	msg := hub.FORWARDIPC{F_IPC: req.Json()}
	ret := _call(packet.Pack(hub.Code["forward_req"], &msg, nil))
	reader := packet.Reader(ret)
	tbl, err := hub.PKT_INT(reader)

	if err != nil || tbl.F_v == 0 {
		return false
	}

	return true
}

//------------------------------------------------ Forward to Hub/Group
func GroupForward(req *IPCObject) bool {
	defer _hub_err()
	msg := hub.FORWARDIPC{F_IPC: req.Json()}
	ret := _call(packet.Pack(hub.Code["forwardgroup_req"], &msg, nil))
	reader := packet.Reader(ret)
	tbl, err := hub.PKT_INT(reader)

	if err != nil || tbl.F_v == 0 {
		return false
	}

	return true
}

func _hub_err() {
	if x := recover(); x != nil {
		log.Println(x)
	}
}
