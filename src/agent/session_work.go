package agent

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
func session_work(sess *Session) bool {
	config := cfg.Get()
	session_timeout := 30 // sec
	if config["session_timeout"] != "" {
		session_timeout, _ = strconv.Atoi(config["session_timeout"])
	}

	if time.Now().Unix()-sess.LastPing > int64(session_timeout) {
		log.Printf("timeout: user %v, connected at: %v\n", sess.Basic.Id, time.Unix(sess.ConnectTime, 0))
		return true
	}

	return false
}
