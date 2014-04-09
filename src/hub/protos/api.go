package protos

import "misc/packet"

var Code = map[string]int16{
	"ping_req":    0,   // PING
	"login_req":   1,   // 登陆
	"logout_req":  2,   // 登出
	"raid_req":    3,   // 攻击
	"protect_req": 4,   // 加保护
	"free_req":    5,   // 结束攻击
	"adduser_req": 6,   // 注册一个新注册的玩家
	"forward_req": 100, // 转发IPC消息
}

var RCode = map[int16]string{
	0:   "ping_req",    // PING
	1:   "login_req",   // 登陆
	2:   "logout_req",  // 登出
	3:   "raid_req",    // 攻击
	4:   "protect_req", // 加保护
	5:   "free_req",    // 结束攻击
	6:   "adduser_req", // 注册一个新注册的玩家
	100: "forward_req", // 转发IPC消息
}

var ProtoHandler map[int16]func(int32, *packet.Packet) []byte

func init() {
	ProtoHandler = map[int16]func(int32, *packet.Packet) []byte{
		0:   P_ping_req,
		1:   P_login_req,
		2:   P_logout_req,
		3:   P_raid_req,
		4:   P_protect_req,
		5:   P_free_req,
		6:   P_adduser_req,
		100: P_forward_req,
	}
}
