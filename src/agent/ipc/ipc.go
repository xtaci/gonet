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
	UNKNOWN = int16(iota)
	SERVICE_PING
)

var IPCHandler map[int16]func(*Session, *IPCObject) bool = map[int16]func(*Session, *IPCObject) bool{
	SERVICE_PING: IPC_ping,
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

	req := &IPCObject{Sender: src_id, Service: service, Object: val, Time: time.Now().Unix()}

	// first try local delivery, if dest_id is not in the same server, just forward to HUB server.
	peer := QueryOnline(dest_id)
	if peer != nil {
		select {
		case peer.MQ <- *req:
		case <-time.After(time.Second):
			panic("deadlock") // rare case, when both chan is full
		}
		return true
	} else {
		// convert req to json again, LEVEL-2 encapsulation
		return _forward(dest_id, req.Json())
	}
}
