package agent

import (
	"time"
	"log"
)

import (
	. "types"
)

func session_work(sess *Session, session_timeout int) bool {
	if time.Now().Unix() - sess.HeartBeat.Unix() > int64(session_timeout) {
		log.Printf("timeout of user %v, occured\n", sess.User.Id)
		return true
	}

	return false
}
