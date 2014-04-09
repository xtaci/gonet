package protos

import (
	"encoding/json"
	"log"
	"runtime"
)

import (
	"db/forward_tbl"
	"helper"
	"hub/core"
	"misc/packet"
	. "types"
)

func HandleRequest(hostid int32, reader *packet.Packet, output chan []byte) {
	defer helper.PrintPanicStack()

	seqid, err := reader.ReadU64() // read seqid
	if err != nil {
		log.Println("Read Sequence Id failed.", err)
		return
	}

	b, err := reader.ReadS16()
	if err != nil {
		log.Println("read protocol error")
		return
	}

	handle := ProtoHandler[b]
	if handle != nil {
		ret := handle(hostid, reader)
		if len(ret) != 0 {
			helper.SendChan(seqid, ret, output)
		}
	}
}

func P_ping_req(hostid int32, reader *packet.Packet) []byte {
	tbl, _ := PKT_INT(reader)
	ret := INT{tbl.F_v}
	return packet.Pack(-1, &ret, nil)
}

func P_login_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_LOGIN_REQ(pkt)
	ret := LOGIN_ACK{F_success: 0}

	if core.Login(tbl.F_id, hostid) {
		ret.F_success = 1
	}

	return packet.Pack(-1, &ret, nil)
}

func P_logout_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if core.Logout(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(-1, &ret, nil)
}

func P_raid_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if core.Raid(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(-1, &ret, nil)
}

func P_protect_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_PROTECT(pkt)
	ret := INT{F_v: 0}

	if core.Protect(tbl.F_id, tbl.F_protecttime) {
		ret.F_v = 1
	}

	return packet.Pack(-1, &ret, nil)
}

func P_free_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if core.Free(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(-1, &ret, nil)
}

func P_forward_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_FORWARDIPC(pkt)

	obj := &IPCObject{}
	err := json.Unmarshal(tbl.F_IPC, obj)

	if err != nil {
		log.Println("decode forward IPCObject error", err)
		return nil
	}

	// to SYS_USR or to player
	if obj.DestID == SYS_USR {
		Syscast(hostid, obj)
	} else {
		_unicast(hostid, obj)
	}

	ret := INT{F_v: 1}
	return packet.Pack(-1, &ret, nil)
}

func _unicast(hostid int32, obj *IPCObject) {
	// if user is online, send to the server, or else send to database
	state := core.State(obj.DestID)

	switch state {
	case ON_PROT, ON_FREE:
		host := core.Host(obj.DestID)
		ch := ForwardChan(host)

		if ch != nil {
			ch <- *obj
		} else {
			forward_tbl.Push(obj)
		}
	default:
		forward_tbl.Push(obj)
	}
}

func Syscast(hostid int32, obj *IPCObject) {
	all := AllServers()
	for _, v := range all {
		if v != hostid { // ignore sender's server
			ch := ForwardChan(v)

			if ch != nil {
				ch <- *obj
			}
		}
	}
}

func P_adduser_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if core.LoadUser(tbl.F_id) {
		core.Login(tbl.F_id, hostid)
		ret.F_v = 1
	}

	return packet.Pack(-1, &ret, nil)
}

func checkErr(err error) {
	if err != nil {
		funcName, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("ERR:%v,[func:%v,file:%v,line:%v]\n", err, runtime.FuncForPC(funcName).Name(), file, line)
		}

		panic("error occured in HUB ipc module")
	}
}
