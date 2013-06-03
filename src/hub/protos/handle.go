package protos

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

import (
	"db/forward_tbl"
	. "helper"
	"hub/accounts"
	"misc/packet"
)

var (
	Servers    map[int32]chan []byte
	ServerLock sync.RWMutex
)

const (
	COLLECTION = "MESSAGES"
)

func init() {
	Servers = make(map[int32]chan []byte)
}

//--------------------------------------------------------- send
func _send(seqid uint64, data []byte, output chan []byte) {
	writer := packet.Writer()
	writer.WriteU16(uint16(len(data)) + 8)
	writer.WriteU64(seqid) // piggyback seq id
	writer.WriteRawBytes(data)
	output <- writer.Data()
}

func HandleRequest(hostid int32, reader *packet.Packet, output chan []byte) {
	defer PrintPanicStack()

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

	fmt.Println("proto: ", b)
	handle := ProtoHandler[b]
	if handle != nil {
		ret := handle(hostid, reader)
		if len(ret) != 0 {
			_send(seqid, ret, output)
		}
	}

}

func P_ping_req(hostid int32, reader *packet.Packet) []byte {
	tbl, _ := PKT_INT(reader)
	ret := INT{tbl.F_v}
	return packet.Pack(Code["ping_ack"], ret, nil)
}

func P_login_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.Login(tbl.F_id, hostid) {
		ret.F_v = 1
	}

	return packet.Pack(Code["login_ack"], ret, nil)
}

func P_logout_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.Logout(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["logout_ack"], ret, nil)
}

func P_changescore_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_CHGSCORE(pkt)
	ret := INT{F_v: 0}

	if accounts.ChangeScore(tbl.F_id, tbl.F_oldscore, tbl.F_newscore) {
		ret.F_v = 1
	}

	return packet.Pack(Code["changescore_ack"], ret, nil)
}

func P_getlist_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_GETLIST(pkt)
	ret := LIST{}

	ids, scores := accounts.GetList(int(tbl.F_A), int(tbl.F_B))
	ret.F_items = make([]ID_SCORE, len(ids))

	for k := range ids {
		ret.F_items[k].F_id = ids[k]
		ret.F_items[k].F_score = scores[k]
	}

	return packet.Pack(Code["getlist_ack"], ret, nil)
}

func P_raid_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.Raid(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["raid_ack"], ret, nil)
}

func P_protect_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_PROTECT(pkt)
	ret := INT{F_v: 0}

	if accounts.Protect(tbl.F_id, tbl.F_protecttime) {
		ret.F_v = 1
	}

	return packet.Pack(Code["protect_ack"], ret, nil)
}

func P_unprotect_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.UnProtect(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["unprotect_ack"], ret, nil)
}

func P_free_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.Free(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["free_ack"], ret, nil)
}

func P_getinfo_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INFO{}
	ret.F_id = tbl.F_id
	ret.F_state = accounts.State(tbl.F_id)
	ret.F_score = accounts.Score(tbl.F_id)
	ret.F_clan = accounts.Score(tbl.F_id)
	ret.F_protecttime = accounts.ProtectTimeout(tbl.F_id)
	ret.F_name = accounts.Name(tbl.F_id)
	if ret.F_state == 0 {
		ret.F_flag = false
	} else {
		ret.F_flag = true
	}

	return packet.Pack(Code["getinfo_ack"], ret, nil)
}

func P_forward_req(hostid int32, pkt *packet.Packet) []byte {
	defer func() {
		if x := recover(); x != nil {
			log.Println("forward packet error")
		}
	}()

	tbl, _ := PKT_FORWARDIPC(pkt)

	// if user is online, send to the server, or else send to database
	state := accounts.State(tbl.F_dest_id)
	host := accounts.Host(tbl.F_dest_id)

	fmt.Println(tbl.F_dest_id, tbl.F_IPC)
	if state&accounts.ONLINE != 0 {
		ServerLock.RLock()
		ch := Servers[host]
		ServerLock.RUnlock()

		ch <- tbl.F_IPC
	} else {
		// send to db
		forward_tbl.Push(tbl.F_dest_id, tbl.F_IPC)
	}

	ret := INT{F_v: 1}
	return packet.Pack(Code["forward_ack"], ret, nil)
}

func P_adduser_req(hostid int32, pkt *packet.Packet) []byte {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.LoadUser(tbl.F_id) {
		accounts.Login(tbl.F_id, hostid)
		ret.F_v = 1
	}

	return packet.Pack(Code["adduser_ack"], ret, nil)
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
