package ipc_service

import (
	"encoding/json"
)

import (
	"agent/net"
	. "helper"
	"misc/packet"
	. "types"
)

//---------------------------------------------------------- 聊天通知
func IPC_chat(sess *Session, obj *IPCObject) []byte {
	w := &Words{}
	err := json.Unmarshal(obj.Object, w)
	if err != nil {
		ERR("cannot decode worldchat message")
		return nil
	}

	ret := &talk{}
	ret.F_user = w.Speaker
	ret.F_msg = w.Words
	return packet.Pack(net.Code["talk_notify"], ret, nil)
}
