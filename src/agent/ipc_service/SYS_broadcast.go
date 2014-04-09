package ipc_service

import (
	"encoding/json"
)

import (
	"agent/gsdb"
	. "helper"
	. "types"
)

//----------------------------------------------------------
// 广播包
// 只由 SYS_USR 接收
// 外部只需要调用Broadcast函数即可
func SYS_broadcast(sess *Session, obj *IPCObject) []byte {
	// 解包
	DEBUG("RECEIVED SYS_BROADCAST")
	realmsg := &IPCObject{}
	err := json.Unmarshal(obj.Object, realmsg)
	if err != nil {
		ERR("SYS_broadcast cannot decode msg", err, obj.Object)
		return nil
	}

	// 投递closure
	send := func(MQ chan IPCObject) {
		defer func() {
			recover()
		}()
		MQ <- *realmsg
	}

	// 循环投递
	DEBUG("SERVICE", realmsg.Service)
	users := gsdb.ListAll()
	for _, v := range users {
		peer := gsdb.QueryOnline(v)
		if peer != nil {
			if v != SYS_USR { // 广播包不能投递给系统玩家
				send(peer.MQ)
			}
		}
	}

	DEBUG("BROADCAST Delivered to", len(users), "users")
	return nil
}
