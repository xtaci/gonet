package main

import (
	"strconv"
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
