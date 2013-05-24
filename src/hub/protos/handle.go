package protos

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"log"
	"runtime"
	"sync"
)

import (
	"cfg"
	. "db"
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
	defer _HandleError()

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
		ret, err := handle(hostid, reader)
		if err == nil {
			_send(seqid, ret, output)
		} else {
			log.Println(ret)
		}
	}

}

func P_ping_req(hostid int32, reader *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_INT(reader)
	ret := INT{tbl.F_v}
	return packet.Pack(Code["ping_ack"], ret, nil), nil
}

func P_login_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.Login(tbl.F_id, hostid) {
		ret.F_v = 1
	}

	return packet.Pack(Code["login_ack"], ret, nil), nil
}

func P_logout_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.Logout(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["logout_ack"], ret, nil), nil
}

func P_changescore_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_CHGSCORE(pkt)
	ret := INT{F_v: 0}

	if accounts.ChangeScore(tbl.F_id, tbl.F_oldscore, tbl.F_newscore) {
		ret.F_v = 1
	}

	return packet.Pack(Code["changescore_ack"], ret, nil), nil
}

func P_getlist_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_GETLIST(pkt)
	ret := LIST{}

	ids, scores := accounts.GetList(int(tbl.F_A), int(tbl.F_B))
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

	if accounts.Raid(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["raid_ack"], ret, nil), nil
}

func P_protect_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.Raid(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["protect_ack"], ret, nil), nil
}

func P_unprotect_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.UnProtect(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["unprotect_ack"], ret, nil), nil
}

func P_free_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.Free(tbl.F_id) {
		ret.F_v = 1
	}

	return packet.Pack(Code["free_ack"], ret, nil), nil
}

func P_getinfo_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INFO{}
	ret.F_id = tbl.F_id
	ret.F_state = accounts.State(tbl.F_id)
	ret.F_score = accounts.Score(tbl.F_id)
	ret.F_clan = accounts.Score(tbl.F_id)
	ret.F_protecttime = accounts.ProtectTimeout(tbl.F_id)
	ret.F_name = accounts.Name(tbl.F_id)
	return packet.Pack(Code["getinfo_ack"], ret, nil), nil
}

func P_forward_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	defer func() {
		if x := recover(); x != nil {
			log.Println("forward packet error")
		}
	}()

	tbl, _ := PKT_FORWARDMSG(pkt)

	// if user is online, send to the server, or else send to database
	state := accounts.State(tbl.F_id)
	host := accounts.Host(tbl.F_id)

	fmt.Println(tbl.F_id, tbl.F_data)
	if state&accounts.ONLINE != 0 {
		ServerLock.RLock()
		ch := Servers[host] // forwarding request
		ServerLock.RUnlock()

		ch <- tbl.F_data
	} else {
		// send to db
		config := cfg.Get()
		c := Mongo.DB(config["mongo_db"]).C(COLLECTION)
		err := c.Insert(tbl)
		if err != nil {
			log.Println(err)
		}
	}

	return nil, nil
}

func P_getofflinemsg_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	// get all forward message from db
	tbl, _ := PKT_ID(pkt)
	var msgs []FORWARDMSG
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)
	err := c.Find(bson.M{"f_id": tbl.F_id}).All(&msgs)
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	ret := OFFLINEMSG{}
	ret.F_msgs = make([]PLAINMSG, len(msgs))

	for k := range msgs {
		ret.F_msgs[k].F_msg = msgs[k].F_data
	}

	// remove messages
	info, err := c.RemoveAll(bson.M{"f_id": tbl.F_id})
	if err != nil {
		log.Println(info, err)
	}

	return packet.Pack(Code["getofflinemsg_ack"], ret, nil), nil
}

func P_adduser_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := PKT_ID(pkt)
	ret := INT{F_v: 0}

	if accounts.LoadUser(tbl.F_id) {
		accounts.Login(tbl.F_id, hostid)
		ret.F_v = 1
	}

	return packet.Pack(Code["adduser_ack"], ret, nil), nil
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
