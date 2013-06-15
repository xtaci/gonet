package main

import (
	"agent/ipc"
	"db/data_tbl"
	"db/user_tbl"
	. "types"
	"types/estates"
	"types/samples"
)

//----------------------------------------------- connection close cleanup work
func close_work(sess *Session) {
	if sess.LoggedIn {
		ipc.Logout(sess.User.Id)
		ipc.UnregisterOnline(sess.User.Id)
		close(sess.MQ)

		// 持久化逻辑#3: 离线时，刷入数据库
		user_tbl.Set(&sess.User)
		data_tbl.Set(estates.COLLECTION, &sess.Estates)
		data_tbl.Set(samples.COLLECTION, &sess.LatencySamples)
	}
}
