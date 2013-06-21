package main

import (
	"log"
)

import (
	"agent/ipc"
	"db/forward_tbl"
	. "types"
)

//----------------------------------------------- connection close cleanup work
func close_work(sess *Session) {
	if sess.LoggedIn {
		ipc.Logout(sess.User.Id)
		ipc.UnregisterOnline(sess.User.Id)
		close(sess.MQ)

		// 未处理的IPC数据，重新放入db
		if len(sess.MQ) > 0 {
			log.Println("re-push ipcobject back to db")
		}

		for len(sess.MQ) > 0 {
			ipcobject := <-sess.MQ
			ipcobject.MarkDelete = false
			forward_tbl.Push(&ipcobject)
		}

		// 持久化逻辑#3: 离线时，刷入数据库
		_flush(sess)
	}
}
