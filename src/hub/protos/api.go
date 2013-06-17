package protos

import "misc/packet"

var Code map[string]int16 = map[string]int16{
	"ping_req":         0,   // PING
	"login_req":        1,   // 登陆
	"logout_req":       2,   // 登出
	"changescore_req":  3,   // 改变分数
	"getlist_req":      4,   // 获取列表
	"raid_req":         5,   // 攻击
	"protect_req":      6,   // 加保护
	"free_req":         7,   // 结束攻击
	"getinfo_req":      8,   // 读取玩家信息
	"adduser_req":      9,   // 注册一个新注册的玩家
	"forward_req":      100, // 转发IPC消息
	"forwardgroup_req": 101, // 转发IPC消息到联盟
}

var RCode map[int16]string = map[int16]string{
	0:   "ping_req",         // PING
	1:   "login_req",        // 登陆
	2:   "logout_req",       // 登出
	3:   "changescore_req",  // 改变分数
	4:   "getlist_req",      // 获取列表
	5:   "raid_req",         // 攻击
	6:   "protect_req",      // 加保护
	7:   "free_req",         // 结束攻击
	8:   "getinfo_req",      // 读取玩家信息
	9:   "adduser_req",      // 注册一个新注册的玩家
	100: "forward_req",      // 转发IPC消息
	101: "forwardgroup_req", // 转发IPC消息到联盟
}

var ProtoHandler map[uint16]func(int32, *packet.Packet) []byte = map[uint16]func(int32, *packet.Packet) []byte{
	0:   P_ping_req,
	1:   P_login_req,
	2:   P_logout_req,
	3:   P_changescore_req,
	4:   P_getlist_req,
	5:   P_raid_req,
	6:   P_protect_req,
	7:   P_free_req,
	8:   P_getinfo_req,
	9:   P_adduser_req,
	100: P_forward_req,
	101: P_forwardgroup_req,
}
