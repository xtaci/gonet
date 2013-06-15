package ipc

import (
	"encoding/json"
)

import (
	. "types"
)

func IPC_ping(sess *Session, obj *IPCObject) bool {
	var str string
	err := json.Unmarshal(obj.Object, &str)
	if err == nil {
		if obj.Sender != -1 {
			Send(-1, obj.Sender, SERVICE_PING, str)
		}
	}
	return true
}
