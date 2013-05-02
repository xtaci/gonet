package ipc

import (
	"errors"
	"hub/names"
	. "types"
)

const (
	UNKNOWN = int16(iota)
	USERINFO_REQUEST
)

type RequestType struct {
	Code   int16            // tos
	CH     chan interface{} // service-oriented data channel
	Params interface{}
}

var RequestHandler map[int16]func(*Session, interface{}) interface{} = map[int16]func(*Session, interface{}) interface{}{
	USERINFO_REQUEST: userinfo_request,
}


//--------------------------------------------------------- Non-Blocking Send
func Send(id int32, tos int16, params interface{}) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("ipc.Send() failed")
		}
	}()

	peer := names.Query(id);
	req := &RequestType{Code: tos}
	req.Params = params
	peer <- req

	return nil
}


//--------------------------------------------------------- Blocking Call
func Call(id int32, tos int16, params interface{}) (ret interface{}, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.New("ipc.Call() failed")
		}
	}()

	peer := names.Query(id);
	req := &RequestType{Code: tos}
	req.CH = make(chan interface{})
	req.Params = params

	peer <- req
	ret = <-req.CH

	return
}
