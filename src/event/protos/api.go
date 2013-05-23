package protos

import "misc/packet"

var Code map[string]int16 = map[string]int16{
	"add_req":    0,    // 添加一个cooldown请求
	"add_ack":    1000, // 返回这个cooldown编号
	"cancel_req": 1,    // 取消一个cooldown请求
	"cancel_ack": 1001, // 返回0
}

var RCode map[int16]string = map[int16]string{
	0:    "add_req",
	1000: "add_ack",
	1:    "cancel_req",
	1001: "cancel_ack",
}

var ProtoHandler map[uint16]func(*packet.Packet) ([]byte, error) = map[uint16]func(*packet.Packet) ([]byte, error){
	0: P_add_req,
	1: P_cancel_req,
}
