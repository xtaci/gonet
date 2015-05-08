package AI

import (
	"db/forward_tbl"
	. "types"
)

//------------------------------------------------ 载入离线时收到的的IPCObject
func LoadIPCObjects(user_id int32, MQ chan IPCObject) {
	objs := forward_tbl.PopAll(user_id)

	// 消息没有完全push到MQ, 存回db
	var k int

	defer func() {
		if x := recover(); x != nil {
			for k < len(objs) {
				forward_tbl.Push(&objs[k])
				k++
			}
		}
	}()

	for k = range objs {
		MQ <- objs[k]
	}
}
