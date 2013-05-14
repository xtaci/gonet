package agent

import (
	"agent/ipc"
	. "types"
)

//----------------------------------------------- connection close cleanup work
func close_work(sess *Session) {
	if sess.IsLoggedOut {		// normal exit
	} else {
	}

	ipc.Unregister(sess.User.Id)
}
