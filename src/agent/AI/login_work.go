package AI

import (
	"encoding/json"
	"log"
	"time"
)

import (
	"db/data_tbl"
	"db/forward_tbl"
	. "types"
)

//------------------------------------------------ 登陆后的数据加载
func LoginWork(sess *Session) bool {
	// 载入建筑表
	data_tbl.Get("ESTATES", sess.User.Id, &sess.Estates)

	// 最后, 载入离线消息，并push到MQ, 这里小心MQ的buffer长度
	objs := forward_tbl.PopAll(sess.User.Id)
	for k := range objs {
		obj := &IPCObject{}
		err := json.Unmarshal(objs[k], obj)
		if err != nil {
			log.Println("illegal IPCObject", objs[k])
		} else {
			sess.MQ <- *obj
		}
	}

	sess.LastFlush = time.Now().Unix()

	return true
}
