package AI

import (
	"encoding/json"
	"log"
)

import (
	"db/forward_tbl"
	. "types"
)

//------------------------------------------------ 载入离线时收到的的IPCObject
func LoadIPCObjects(user_id int32, MQ chan IPCObject) {
	objs := forward_tbl.PopAll(user_id)
	for k := range objs {
		obj := &IPCObject{}
		err := json.Unmarshal(objs[k], obj)
		if err != nil {
			log.Println("Illegal IPCObject", objs[k])
		} else {
			MQ <- *obj
		}
	}
}
