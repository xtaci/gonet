package ipc

import (
	. "types"
)

func userinfo_request(sess *Session, params interface{}) (interface{}, []byte) {
	return sess.User, nil
}
