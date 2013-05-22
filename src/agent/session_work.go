package agent

import (
	"log"
	"time"
)

import (
	. "types"
)

//----------------------------------------------- session timeout
func session_work(sess *Session, session_timeout int) bool {
	if time.Now().Unix()-sess.HeartBeat.Unix() > int64(session_timeout) {
		log.Printf("timeout of user %v, occured\n", sess.Basic.Id)
		return true
	}

	return false
}
