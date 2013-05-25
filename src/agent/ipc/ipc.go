package ipc

import (
	"encoding/json"
	"log"
	"time"
)

import (
	. "types"
)

const (
	IPC_PING = int16(1)
)

var IPCHandler map[int16]func(*Session, *IPCObject) = map[int16]func(*Session, *IPCObject){
	IPC_PING: IPC_ping,
}

//------------------------------------------------ p2p send from src_id to dest_id
func Send(src_id, dest_id int32, service int16, object interface{}) (ret bool) {
	defer func() {
		if x := recover(); x != nil {
			ret = false
		}
	}()

	// convert the OBJECT to json, LEVEL-1 encapsulation
	val, err := json.Marshal(object)
	if err != nil {
		log.Println("IPC Send error:", err)
		return false
	}
	req := IPCObject{Sender: src_id, Service: service, Object: val}

	// first try local delivery, if dest_id is not in same server, forward to hub
	peer := QueryOnline(dest_id)
	if peer != nil {
		select {
		case peer.MQ <- req:
		case <-time.After(time.Second):
			panic("deadlock") // rare case, when both chan is full
		}
		return
	} else {
		// convert req to json again, LEVEL-2 encapsulation
		req_json, _ := json.Marshal(object)
		return ForwardHub(dest_id, req_json)
	}

	return true
}
