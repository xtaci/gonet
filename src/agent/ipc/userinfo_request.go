package ipc

import (
	. "types"
)

func userinfo_request(sess *Session, request *RequestType)(ret []byte, err error) {
	request.CH <- sess.User
	return nil,nil
}
