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

//----------------------------------------------- session timeout
func session_timeout(sess *Session) bool {
	config := cfg.Get()
	timeout := 30 // sec
	if config["session_timeout"] != "" {
		timeout, _ = strconv.Atoi(config["session_timeout"])
	}

	if time.Now().Unix()-sess.LastPacketTime > int64(timeout) {
		log.Printf("timeout: user %v, connected at: %v\n", sess.User.Id, sess.ConnectTime)
		return true
	}

	return false
}
