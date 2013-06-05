package protos

import "misc/packet"

var Code map[string]int16 = map[string]int16{
	"ping_req": 0,    // PING
	"ping_ack": 1000, // 返回请求数值
	"add_req":  1,    // 加入一个统计数
	"add_ack":  1001, // 返回0
}

var RCode map[int16]string = map[int16]string{
	0:    "ping_req", // PING
	1000: "ping_ack", // 返回请求数值
	1:    "add_req",  // 加入一个统计数
	1001: "add_ack",  // 返回0
}

var ProtoHandler map[uint16]func(*packet.Packet) []byte = map[uint16]func(*packet.Packet) []byte{
	0: P_ping_req,
	1: P_add_req,
}
