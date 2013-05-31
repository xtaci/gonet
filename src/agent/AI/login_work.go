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
	"types/estates"
	"types/grid"
	"misc/naming"
)

//------------------------------------------------ 登陆后的数据加载
func LoginWork(sess *Session) bool {
	// 载入建筑表
	data_tbl.Get(estates.COLLECTION, sess.User.Id, &sess.Estates)
	// 建立位图的格子信息
	sess.Grid = grid.New()
	for _, v := range sess.Estates.Estates {
		// TODO :  读gamedata,建立grid信息
		sess.Grid.Set(v.X, v.Y, naming.FNV1a(v.TYPE))
	}

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

	sess.LastFlushTime = time.Now().Unix()

	return true
}
