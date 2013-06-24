package ipc

import (
	"encoding/json"
)

import (
	"misc/packet"
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
		ret := &talk{F_user: sess.User.Name, F_msg: str}
		return packet.Pack(Code["talk_notify"], ret, nil)
	}

	return nil
}
