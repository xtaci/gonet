package main

import (
	"strconv"
	"time"
)

import (
	"cfg"
	"db/user_tbl"
	"helper"
	. "types"
)

//------------------------------------------------ data flush control (interval + dirty flag)
func _flush_work(sess *Session) {
	config := cfg.Get()
	fi := config["flush_interval"]
	inter, _ := strconv.Atoi(fi)
	fo := config["flush_ops"]
	ops, _ := strconv.Atoi(fo)

	if sess.DirtyCount() > int32(ops) || (sess.DirtyCount() > 0 && time.Now().Unix()-sess.User.LastSaveTime > int64(inter)) {
		helper.NOTICE("flush dirtycount:", sess.DirtyCount(), "duration(sec):", time.Now().Unix()-sess.User.LastSaveTime)
		_flush(sess)
	}
}

//------------------------------------------------ save to db
func _flush(sess *Session) {
	if sess.User != nil {
		sess.User.LastSaveTime = time.Now().Unix()
		user_tbl.Set(sess.User)
		helper.NOTICE(sess.User.Id, sess.User.Name, "data flushed")
	}

	// TODO : save all the data in session
	sess.MarkClean()
}
