package AI

import (
	"encoding/json"
	"log"
)

import (
	"db/estate_tbl"
	"db/forward_tbl"
	. "types"
)

//------------------------------------------------ 登陆后的数据加载
func LoginWork(sess *Session) bool {
	// 从数据库中载入玩家数据
	sess.EstateManager = estate_tbl.Get(sess.Basic.Id)
	// 载入建筑表
	// 载入离线消息，并push到MQ, 这里小心MQ的buffer长度
	objs := forward_tbl.PopAll(sess.Basic.Id)
	for k := range objs {
		obj := &IPCObject{}
		err := json.Unmarshal(objs[k], obj)
		if err != nil {
			log.Println("illegal IPCObject", objs[k])
		} else {
			sess.MQ <- *obj
		}
	}

	return true
}
