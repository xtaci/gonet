package types

import (
	"encoding/json"
)

//---------------------------------------------------------- IPCObject 定义
type IPCObject struct {
	SrcID      int32   // 发送方用户ID
	DestID     int32   // 接收放用户ID
	AuxIDs     []int32 `bson:",omitempty"` // 目标用户ID集合(用于组播)
	Service    int16   // 服务号
	Object     []byte  // 投递的 JSON STRING
	Time       int64   // 发送时间
	MarkDelete bool    // 数据库标记删除
}

//---------------------------------------------------------- 将整个IPCObject转为JSON
func (obj *IPCObject) Json() []byte {
	val, _ := json.Marshal(obj)
	return val
}
