package ipc_service

import (
	. "helper"
	. "types"
)

//---------------------------------------------------------- 用于强制挤下线
func IPC_kick(sess *Session, obj *IPCObject) []byte {
	NOTICE(sess.User.Name, "Kicked Out")
	sess.Flag |= SESS_KICKED_OUT
	return nil
}
