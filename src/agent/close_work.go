package agent

import (
	"agent/ipc"
	. "types"
)

//----------------------------------------------- connection close cleanup work
func close_work(sess *Session) {
	if sess.LoggedIn {
		// TODO: flush db 
		ipc.UnregisterOnline(sess.Basic.Id)
	}
}
