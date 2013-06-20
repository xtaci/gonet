package protos

import "misc/packet"

var Code map[string]int16 = map[string]int16{
	"ping_req":   0,    // PING
	"ping_ack":   1000, // 返回请求数值
	"add_req":    1,    // 添加一个离线事件
	"add_ack":    1001, // 返回事件编号
	"cancel_req": 2,    // 取消一个事件
	"cancel_ack": 1002, // 返回0
}

var RCode map[int16]string = map[int16]string{
	0:    "ping_req",   // PING
	1000: "ping_ack",   // 返回请求数值
	1:    "add_req",    // 添加一个离线事件
	1001: "add_ack",    // 返回事件编号
	2:    "cancel_req", // 取消一个事件
	1002: "cancel_ack", // 返回0
}

var ProtoHandler map[uint16]func(*packet.Packet) []byte = map[uint16]func(*packet.Packet) []byte{
	0: P_ping_req,
	1: P_add_req,
	2: P_cancel_req,
}
