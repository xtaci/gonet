package ipc

import (
	. "types"
)

func userinfo_request(sess *Session, params interface{}) (ret interface{}) {
	return sess.User
}
