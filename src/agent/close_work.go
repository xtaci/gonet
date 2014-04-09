package main

import (
	"agent/gsdb"
	"agent/hub_client"
	"db/forward_tbl"
	. "helper"
	"misc/geoip"
	. "types"
)

//----------------------------------------------- 会话结束时的扫尾清理工作
func close_work(sess *Session) {
	defer PrintPanicStack()
	if sess.Flag&SESS_LOGGED_IN == 0 {
		return
	}

	// 退出时玩家数据<<强制>>刷入数据库
	_flush(sess)

	// 通知HUB离线
	hub_client.Logout(sess.User.Id)
	// 在GSDB反注册
	gsdb.UnregisterOnline(sess.User.Id)
	// 关闭chan时，如果有sender，有可能panic，注意IPC对panic的处理
	close(sess.MQ)
	// 未处理的IPC数据，重新放入db, 注意:必须先close MQ
	for ipcobject := range sess.MQ {
		forward_tbl.Push(&ipcobject)
		NOTICE("re-pushed ipcobject back to db, userid:", sess.User.Id)
	}

	NOTICE(sess.User.Name, "disconnected from", sess.IP, "country:", geoip.Query(sess.IP))
}
