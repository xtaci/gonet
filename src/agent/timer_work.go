package main

import (
	"strconv"
	"time"
)

import (
	"cfg"
	. "types"
)

//----------------------------------------------- timer work
func timer_work(sess *Session) {
	// check whether the user is logged in
	if !sess.LoggedIn {
		return
	}

	// TODO: 持久化逻辑#2： 超过一定的时间，刷入数据库
	config := cfg.Get()
	ivl, _ := strconv.Atoi(config["flush_interval"])
	if sess.Dirty && time.Now().Unix()-sess.LastFlushTime > int64(ivl) {
		// 刷入数据到数据库
		_flush(sess)
	}
}
