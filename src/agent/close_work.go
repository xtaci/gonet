package main

import (
	"log"
)

import (
	"agent/gsdb"
	"agent/hub_client"
	"db/forward_tbl"
	. "types"
)

//----------------------------------------------- connection close cleanup work
func close_work(sess *Session) {
	if sess.LoggedIn {
		hub_client.Logout(sess.User.Id)
		gsdb.UnregisterOnline(sess.User.Id)
		close(sess.MQ)

		// un-delivered IPCObject(s) should be saved back to database!
		for ipcobject := range sess.MQ {
			forward_tbl.Push(&ipcobject)
			log.Println("re-push ipcobject back to db")
		}

		// Data Persistence #3: offline flush
		_flush(sess)
	}
}
