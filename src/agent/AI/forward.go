package AI

import (
	"db/forward_tbl"
	. "types"
)

//------------------------------------------------ 载入离线时收到的的IPCObject
func LoadIPCObjects(user_id int32, MQ chan IPCObject) {
	objs := forward_tbl.PopAll(user_id)
	for k := range objs {
		MQ <- objs[k]
	}
}
