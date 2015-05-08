package main

import (
	"agent/gsdb"
	"agent/hub_client"
	"db/forward_tbl"
	. "helper"
	"misc/geoip"
	. "types"
)

//----------------------------------------------- cleanup work after disconnection
func close_work(sess *Session) {
	defer PrintPanicStack()
	if sess.Flag&SESS_LOGGED_IN == 0 {
		return
	}

	// must flush user data
	_flush(sess)

	// notify hub
	hub_client.Logout(sess.User.Id)

	// unregister online at this server
	gsdb.UnregisterOnline(sess.User.Id)

	// close MQ, and save the queue to db
	close(sess.MQ)
	for ipcobject := range sess.MQ {
		forward_tbl.Push(&ipcobject)
		NOTICE("re-pushed ipcobject back to db, userid:", sess.User.Id)
	}

	NOTICE(sess.User.Name, "disconnected from", sess.IP, "country:", geoip.Query(sess.IP))
}
