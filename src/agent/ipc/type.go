package ipc

import (
	. "types"
	"hub/names"
	"errors"
)

const (
	UNKNOWN = int16(iota)
	USERINFO_REQUEST
)

type RequestType struct {
	Code int16			// tos
	CH chan interface{} // service-oriented data channel
	Params interface{}
}

var RequestHandler map[int16]func(*Session, *RequestType) ([]byte, error) = map[int16]func(*Session, *RequestType)([]byte, error){
	USERINFO_REQUEST:userinfo_request,
}

func Call(id int32, tos int16, params interface{}) (ret interface{}, err error) {
	if peer := names.Query(id);peer!=nil {
		req := &RequestType{Code:tos}
		req.CH = make(chan interface{}, 1)
		req.Params = params
		peer <- req

		ret, ok := <-req.CH

		if ok {
			return ret, nil
		}
	}

	return nil, errors.New("Call() failed")
}
