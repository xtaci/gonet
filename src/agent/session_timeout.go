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
	DEFAULT_TIMEOUT = 30
)

//----------------------------------------------- session timeout
func session_timeout(sess *Session) bool {
	config := cfg.Get()
	timeout, err := strconv.Atoi(config["session_timeout"])
	if err != nil {
		log.Println("cannot parse session_timeout from config", err)
		timeout = DEFAULT_TIMEOUT
	}

	// compare current with last packet arrival time
	if time.Now().Unix()-sess.LastPacketTime > int64(timeout) {
		if sess.LoggedIn {
			log.Printf("timeout: user %v, connected at: %v\n", sess.User.Id, sess.ConnectTime)
		}
		return true
	}

	return false
}
