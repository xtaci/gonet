package main

import (
	"log"
	"strconv"
	"time"
)

import (
	"cfg"
	. "types"
)

const (
	DEFAULT_FLUSH_INTERVAL = 30
)

//----------------------------------------------- timer work
func timer_work(sess *Session) {
	// check whether the user is logged in
	if !sess.LoggedIn {
		return
	}

	// 持久化逻辑#2： 超过一定的时间，刷入数据库
	config := cfg.Get()
	ivl, err := strconv.Atoi(config["flush_interval"])
	if err != nil {
		log.Println("cannot parse flush_interval from config", err)
		ivl = DEFAULT_FLUSH_INTERVAL
	}

	if sess.OpCount > 0 && time.Now().Unix()-sess.LastFlushTime > int64(ivl) {
		// 刷入数据到数据库
		_flush(sess)
	}
}
