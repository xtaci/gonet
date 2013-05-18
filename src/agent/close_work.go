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

	ipc.UnregisterOnline(sess.User.Id)
}
