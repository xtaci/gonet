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
		if obj.SrcID != -1 {
			Send(-1, obj.SrcID, SERVICE_PING, false, str)
		}
	}
	return true
}
