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

	b, err := reader.ReadU16()
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
	ret := LOGIN_ACK{F_success: false}

	if core.Login(tbl.F_id, hostid) {
		ret.F_success = true

		// 登陆后，将联盟消息push给玩家
		group := core.Group(tbl.F_group)
		if group != nil {
			ch := ForwardChan(hostid)
			objs := group.Recv(tbl.F_groupmsgmax + 1)
			for _, v := range objs {
				ch <- v.Json()
			}

			ret.F_groupmsgmax = group.MaxMsgId
		}
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

func P_changescore_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_CHGSCORE(pkt)
	ret := INT{F_v: 0}

	if core.UpdateScore(tbl.F_id, tbl.F_oldscore, tbl.F_newscore) {
		ret.F_v = 1
	}

	return packet.Pack(-1, &ret, nil)
}

func P_getlist_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_GETLIST(pkt)
	ret := LIST{}

	ids, scores := core.GetList(int(tbl.F_A), int(tbl.F_B))
	ret.F_items = make([]ID_SCORE, len(ids))

	for k := range ids {
		ret.F_items[k].F_id = ids[k]
		ret.F_items[k].F_score = scores[k]
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

func P_getinfo_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INFO{}
	ret.F_id = tbl.F_id
	ret.F_state = core.State(tbl.F_id)
	ret.F_score = core.Score(tbl.F_id)
	ret.F_protecttime = core.ProtectTimeout(tbl.F_id)
	if ret.F_state == 0 {
		ret.F_flag = false
	} else {
		ret.F_flag = true
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

	switch obj.CastType {
	case UNICAST:
		_unicast(hostid, obj)
	case MULTICAST:
		_multicast(hostid, obj)
	case GLOBAL_BROADCAST:
		_broadcast(hostid, obj)
	default:
		log.Println("CastType error", hostid, obj)
	}

	ret := INT{F_v: 1}
	return packet.Pack(-1, &ret, nil)
}

func _unicast(hostid int32, obj *IPCObject) {
	// if user is online, send to the server, or else send to database
	state := core.State(obj.DestID)

	switch state {
	case core.ON_PROT, core.ON_FREE:
		host := core.Host(obj.DestID)

		ch := ForwardChan(host)

		if ch != nil {
			ch <- obj.Json()
		} else {
			forward_tbl.Push(obj)
		}
	default:
		forward_tbl.Push(obj)
	}
}

func _broadcast(hostid int32, obj *IPCObject) {
	all := AllServers()
	for _, v := range all {
		if v != hostid { // ignore sender server
			host := core.Host(v)
			ch := ForwardChan(host)

			if ch != nil {
				ch <- obj.Json()
			}
		}
	}
}

func _multicast(hostid int32, obj *IPCObject) {
	group := core.Group(obj.DestID)
	if group == nil {
		log.Println("forward ipc: no such group")
		return
	}
	group.Push(obj)

	// send to online users directly
	members := group.Members()
	for _, user_id := range members {
		state := core.State(user_id)
		switch state {
		case core.ON_PROT, core.ON_FREE:
			host := core.Host(user_id)

			ch := ForwardChan(host)
			if ch != nil {
				ch <- obj.Json()
			}
		}
	}

	return
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
