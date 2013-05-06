package ipc

import (
	"errors"
	"time"
)

import (
	"hub/online"
	. "types"
)

const (
	UNKNOWN = int16(iota)
	USERINFO_REQUEST
	PLAIN_DATA				// userful for server generated packet
)

type RequestType struct {
	Code   int16            // tos
	CH     chan interface{} // service-oriented data channel
	Params interface{}
}

//--------------------------------------------------------- return to ipc && bytes to client
var RequestHandler map[int16]func(*Session, interface{})(interface{}, []byte) = map[int16]func(*Session, interface{})(interface{}, []byte){
	USERINFO_REQUEST: userinfo_request,
}

//--------------------------------------------------------- Non-Blocking Send
func Send(id int32, tos int16, params interface{}) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("ipc.Send() failed")
		}
	}()

	peer := online.Query(id)
	req := &RequestType{Code: tos}
	req.Params = params
	select {
	case peer.MQ <- req:
	case <-time.After(time.Second):
		panic("deadlock") // rare case, when both chan is full
	}

	return nil
}

//--------------------------------------------------------- Blocking Call
func Call(id int32, tos int16, params interface{}) (ret interface{}, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("ipc.Call() failed")
		}
	}()

	peer := online.Query(id)
	req := &RequestType{Code: tos}
	req.CH = make(chan interface{})
	req.Params = params

	select {
	case peer.CALL <- req: // panic on closed channel
		ret = <-req.CH
	case <-time.After(time.Second):
		panic("deadlock")
	}

	return
}
