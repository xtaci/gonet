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
	// if the user is not logged in
	// just return
	if !sess.LoggedIn {
		return
	}

	// Data Persistence #2：Exceeds max flush interval? flush!
	config := cfg.Get()
	ivl, err := strconv.Atoi(config["flush_interval"])
	if err != nil {
		log.Println("cannot parse flush_interval from config", err)
		ivl = DEFAULT_FLUSH_INTERVAL
	}

	if sess.OpCount > 0 && time.Now().Unix()-sess.LastFlushTime > int64(ivl) {
		_flush(sess)
	}
}
