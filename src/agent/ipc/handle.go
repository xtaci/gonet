package ipc

import (
	"encoding/json"
)

import (
	. "types"
)

func IPC_ping(sess *Session, obj *IPCObject) []byte {
	var str string
	err := json.Unmarshal(obj.Object, &str)
	if err == nil {
		if obj.SrcID != -1 {
			Send(-1, obj.SrcID, SERVICE_PING, false, str)
		}
	}
	return nil
}

func IPC_talk(sess *Session, obj *IPCObject) []byte {
	var str string
	err := json.Unmarshal(obj.Object, &str)
	if err == nil {
		return []byte(str)
	}

	return nil
}
