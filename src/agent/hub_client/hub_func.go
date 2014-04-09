package hub_client

import (
	. "helper"
	"misc/packet"
	. "types"
)

type Info struct {
	Id          int32
	State       byte
	Score       int32
	Rank        int32
	ProtectTime int64
}

//---------------------------------------------------------- 测试HUB连通性
func Ping() bool {
	defer _hub_err()
	req := INT{}
	req.F_v = 1
	ret := _call(packet.Pack(Code["ping_req"], &req, nil))
	if ret == nil {
		ERR("HUB Ping return", false)
		return false
	}
	return true
}

//---------------------------------------------------------- 读取某玩家状态
func GetState(id int32) int32 {
	defer _hub_err()

	req := ID{}
	req.F_id = id
	ret := _call(packet.Pack(Code["getstate_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, _ := PKT_INT(reader)

	return tbl.F_v
}

//---------------------------------------------------------- 新用户注册，告知HUB
func AddUser(id int32) bool {
	defer _hub_err()
	req := ID{}
	req.F_id = id
	ret := _call(packet.Pack(Code["adduser_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, err := PKT_INT(reader)

	if err != nil || tbl.F_v == 0 {
		NOTICE("HUB AddUser  return", false, "param", id)
		return false
	}

	return true
}

//---------------------------------------------------------- FSM:  登陆
func Login(id int32) bool {
	defer _hub_err()
	req := LOGIN_REQ{}
	req.F_id = id
	ret := _call(packet.Pack(Code["login_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, _ := PKT_LOGIN_ACK(reader)

	if tbl.F_success == 0 {
		NOTICE("HUB Login return  ", false, "param", id)
		return false
	}

	return true
}

//---------------------------------------------------------- FSM:  登出
func Logout(id int32) bool {
	defer _hub_err()
	req := ID{}
	req.F_id = id
	ret := _call(packet.Pack(Code["logout_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, _ := PKT_INT(reader)

	if tbl.F_v == 0 {
		NOTICE("HUB Logout return  ", false, "param", id)
		return false
	}

	return true
}

//---------------------------------------------------------- FSM:  攻击
func Raid(id int32) bool {
	defer _hub_err()

	req := ID{}
	req.F_id = id
	ret := _call(packet.Pack(Code["raid_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, _ := PKT_INT(reader)

	if tbl.F_v == 0 {
		NOTICE("HUB Raid return ", false, "param", id)
		return false
	}

	return true
}

//---------------------------------------------------------- FSM: 保护
func Protect(id int32, until int64) bool {
	defer _hub_err()

	req := PROTECT{}
	req.F_id = id
	req.F_protecttime = until
	ret := _call(packet.Pack(Code["protect_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, _ := PKT_INT(reader)

	if tbl.F_v == 0 {
		NOTICE("HUB Protect return", false, "param", id, until)
		return false
	}

	return true
}

//---------------------------------------------------------- FSM: 撤销保护
func Free(id int32) bool {
	defer _hub_err()

	req := ID{}
	req.F_id = id
	ret := _call(packet.Pack(Code["free_req"], &req, nil))
	reader := packet.Reader(ret)
	tbl, _ := PKT_INT(reader)

	if tbl.F_v == 0 {
		NOTICE("HUB Free return", false, "param", id)
		return false
	}

	return true
}

//---------------------------------------------------------- 转发IPCObject
func Forward(req *IPCObject) bool {
	defer _hub_err()
	// HUB protocol forwarding
	msg := FORWARDIPC{F_IPC: req.Json()}
	ret := _call(packet.Pack(Code["forward_req"], &msg, nil))
	reader := packet.Reader(ret)
	tbl, err := PKT_INT(reader)

	if err != nil || tbl.F_v == 0 {
		ERR("HUB Forward return false", "param", req)
		return false
	}

	return true
}

func _hub_err() {
	if x := recover(); x != nil {
		ERR(x)
	}
}
