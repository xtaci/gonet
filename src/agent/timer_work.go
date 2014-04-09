package main

import (
	"strconv"
	"sync/atomic"
	"time"
)

import (
	"cfg"
	"helper"
	. "types"
)

//----------------------------------------------- user's timer
func timer_work(sess *Session) {
	if sess.Flag&SESS_LOGGED_IN == 0 {
		return
	}

	// SIGTERM check
	if atomic.LoadInt32(&SIGTERM) == 1 {
		sess.Flag |= SESS_KICKED_OUT
		helper.NOTICE("SIGTERM received, user exits.", sess.User.Id, sess.User.Name)
	}

	// limit rate of request per minute
	config := cfg.Get()
	rpm_limit, _ := strconv.ParseFloat(config["rpm_limit"], 32)
	rpm := float64(sess.PacketCount) / float64(time.Now().Unix()-sess.ConnectTime.Unix()) * 60

	if rpm > rpm_limit {
		sess.Flag |= SESS_KICKED_OUT
		helper.ERR("user RPM too high", sess.User.Id, sess.User.Name, "RPM:", rpm)
		return
	}

	// try save the data
	_flush_work(sess)
}
