package ipc_service

import (
	"encoding/json"
)

import (
	. "agent/ipc"
	. "types"
)

//---------------------------------------------------------- 连通性测试
func IPC_ping(sess *Session, obj *IPCObject) []byte {
	var str string
	err := json.Unmarshal(obj.Object, &str)
	if err == nil {
		if obj.SrcID != 0 {
			Send(0, obj.SrcID, SVC_PING, str)
		}
	}
	return nil
}
