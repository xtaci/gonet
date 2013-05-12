package ipc

import (
	"errors"
	"time"
)

import (
	"agent/online"
	"misc/packet"
	. "types"
)

const (
	UNKNOWN = int16(iota)
	ECHO
)

type RequestType struct {
	Sender int32 // player id
	Code   int16
	Data   []byte
}

//--------------------------------------------------------- return to ipc && bytes to client
var RequestHandler map[int16]func(*Session, []byte) []byte = map[int16]func(*Session, []byte) []byte{}

//--------------------------------------------------------- Non-Blocking Send
func Send(id int32, tos int16, data []byte) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("ipc.Send() failed")
		}
	}()

	peer := online.Query(id)
	req := &RequestType{Code: tos}
	req.Data = data

	if peer != nil { // local delivery
		select {
		case peer.MQ <- req:
		case <-time.After(time.Second):
			panic("deadlock") // rare case, when both chan is full
		}
		return
	} else { // delivery to hub
		return ForwardHub(id, packet.Pack(-1, req, nil))
	}

	return
}
