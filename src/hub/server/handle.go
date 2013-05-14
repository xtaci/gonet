package protos

import (
	"hub/ranklist"
	"misc/packet"
	"github.com/hoisie/redis"
)

import (
	"fmt"
	"log"
	"runtime"
)

var _redis redis.Client

func init() {
	_redis.Addr = "127.0.0.1:6379"
}

func HandleRequest(hostid int32, reader *packet.Packet) []byte {
	defer _HandleError()

	b, err := reader.ReadU16()

	if err != nil {
		log.Println("read protocol error")
	}

	fmt.Println("proto: ", b)
	handle := ProtoHandler[b]
	if handle != nil {
		ret, err := handle(hostid, reader)
		fmt.Println("ret:", ret)
		if err == nil {
			return ret
		}
	}

	return nil
}

func P_forward_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	defer func() {
		if x := recover(); x != nil {
			log.Println("forward packet error")
		}
	}()

	tbl, _ := PKT_MSG(pkt)

	// if user is online, send to the user, or else send to redis 
	state := ranklist.State(tbl.F_id)
	host := ranklist.Host(tbl.F_id)

	fmt.Println(tbl.F_id, tbl.F_data)
	if state&ranklist.ONLINE != 0 {
		_server_lock.RLock()
		ch := _servers[host] // forwarding request
		_server_lock.RUnlock()

		ch <- tbl.F_data
	} else {
		// send to redis
		_redis.Rpush(fmt.Sprintf("MSG#%v", tbl.F_id), tbl.F_data)
	}

	return nil, nil
}

func P_login_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if ranklist.Login(tbl.F_id, hostid) {
		ret.F_v = 1
	}

	return packet.Pack(Code["login_ack"], ret, nil), nil
}

func P_logout_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if ranklist.Logout(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["logout_ack"], ret, nil), nil
}

func P_changescore_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_CHGSCORE(pkt)
	ret := INT{F_v: 0}

	if ranklist.ChangeScore(tbl.F_id, tbl.F_oldscore, tbl.F_newscore) {
		ret.F_v = 1
	}

	return packet.Pack(Code["changescore_ack"], ret, nil), nil
}

func P_getlist_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_GETLIST(pkt)
	ret := LIST{}

	ids, scores := ranklist.GetList(int(tbl.F_A), int(tbl.F_B))
	ret.F_items = make([]ID_SCORE, len(ids))

	for k := range ids {
		ret.F_items[k].F_id = ids[k]
		ret.F_items[k].F_score = scores[k]
	}

	return packet.Pack(Code["getlist_ack"], ret, nil), nil
}

func P_raid_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if ranklist.Raid(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["raid_ack"], ret, nil), nil
}

func P_protect_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if ranklist.Raid(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["protect_ack"], ret, nil), nil
}

func P_unprotect_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if ranklist.Unprotect(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["unprotect_ack"], ret, nil), nil
}

func P_free_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if ranklist.Free(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["free_ack"], ret, nil), nil
}

func P_getinfo_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INFO{}
	ret.F_id = tbl.F_id
	ret.F_state = ranklist.State(tbl.F_id)
	ret.F_score = ranklist.Score(tbl.F_id)
	ret.F_clan = ranklist.Score(tbl.F_id)
	ret.F_protecttime = ranklist.ProtectTime(tbl.F_id)
	ret.F_name = ranklist.Name(tbl.F_id)
	return packet.Pack(Code["getinfo_ack"], ret, nil), nil
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

func _HandleError() {
	if x := recover(); x != nil {
		log.Printf("run time panic when processing HUB request: %v", x)
		for i := 0; i < 10; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				log.Printf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			}
		}
	}
}
