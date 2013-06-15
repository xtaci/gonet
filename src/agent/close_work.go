package main

import (
	"agent/ipc"
	. "types"
)

//----------------------------------------------- connection close cleanup work
func close_work(sess *Session) {
	if sess.LoggedIn {
		// TODO: 持久化逻辑#3: 离线时，刷入数据库
		ipc.Logout(sess.User.Id)
		ipc.UnregisterOnline(sess.User.Id)
	}
}
