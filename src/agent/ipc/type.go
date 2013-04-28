package ipc

import (
	. "types"
)

const (
	UNKNOWN = iota
	USERINFO_REQUEST
)

type RequestType struct {
	Code int16			// tos
	CH chan interface{} // service-oriented data channel
}

var RequestHandler map[int16]func(*Session, *RequestType) ([]byte, error) = map[int16]func(*Session, *RequestType)([]byte, error){
	USERINFO_REQUEST:userinfo_request,
}
