package ipc

import (
	"encoding/json"
	"log"
	"time"
)

import (
	"agent/gsdb"
	"agent/hub_client"
	"db/forward_tbl"
	. "types"
)

const (
	UNKNOWN = int16(iota)
	SERVICE_PING
	SERVICE_TALK
)

var IPCHandler map[int16]func(*Session, *IPCObject) []byte = map[int16]func(*Session, *IPCObject) []byte{
	SERVICE_PING: IPC_ping,
	SERVICE_TALK: IPC_talk,
}

//---------------------------------------------------------- p2p send from src_id to dest_id
func Send(src_id, dest_id int32, service int16, casttype int32, object interface{}) (ret bool) {
	// convert the OBJECT to json, LEVEL-1 encapsulation
	val, err := json.Marshal(object)
	if err != nil {
		log.Println("cannot marshal object to json", err)
		return false
	}

	req := &IPCObject{SrcID: src_id,
		DestID:   dest_id,
		Service:  service,
		CastType: casttype,
		Object:   val,
		Time:     time.Now().Unix()}

	switch casttype {
	case MULTICAST:
		hub_client.Forward(req)
	case GLOBAL_BROADCAST:
		hub_client.Forward(req)
		fallthrough
	case LOCAL_BROADCAST:
		users := gsdb.ListAll()
		defer func() {
			recover()
		}()

		for _, v := range users {
			peer := gsdb.QueryOnline(v)
			if peer != nil {
				select {
				case peer.MQ <- *req:
				case <-time.After(time.Second):
					panic("deadlock") // rare case, when both chan is full
				}
			}
		}
		return true

	case UNICAST:
		// first try local delivery, if dest_id is not in the same server, just forward to HUB server.
		peer := gsdb.QueryOnline(dest_id)
		if peer != nil {
			defer func() {
				if x := recover(); x != nil {
					ret = false
					forward_tbl.Push(req)
				}
			}()

			select {
			case peer.MQ <- *req:
			case <-time.After(time.Second):
				panic("deadlock") // rare case, when both chan is full
			}
			return true
		} else {
			// convert req to json again, LEVEL-2 encapsulation
			return hub_client.Forward(req)
		}
	}

	return false
}
