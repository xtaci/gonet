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
// 组播包
func SYS_multicast(sess *Session, obj *IPCObject) []byte {
	DEBUG("RECEIVED SYS_MULTICAST")
	realmsg := &IPCObject{}
	err := json.Unmarshal(obj.Object, realmsg)
	if err != nil {
		ERR("SYS_multicast cannot decode msg", err, obj.Object)
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
	DEBUG("SERVICE:", realmsg.Service)
	for _, v := range obj.AuxIDs {
		peer := gsdb.QueryOnline(v)
		if peer != nil {
			// if target userid is SYS_USR, it's possible that SYS_USR's MQ is full,
			// and deadlock will happen, so, we use GO!
			if v == SYS_USR {
				go send(peer.MQ)
			} else {
				send(peer.MQ)
			}
		}
	}

	DEBUG("MULTICAST Delivered to", len(obj.AuxIDs), "users")
	return nil
}
