package protos

import "misc/packet"

var Code map[string]int16 = map[string]int16 {
	"forward_req":0,	// 消息转发
	"forward_ack":1000,	// 
	"login_req":1,	// 登陆
	"login_ack":1001,	// 
	"logout_req":2,	// 登出
	"logout_ack":1002,	// 
	"changescore_req":3,	// 改变分数
	"changescore_ack":1003,	// 
	"getlist_req":4,	// 获取列表
	"getlist_ack":1004,	// 
	"raid_req":5,	// 攻击
	"raid_ack":1005,	// 
	"protect_req":6,	// 加保护
	"protect_ack":1006,	// 
	"unprotect_req":7,	// 撤销保护
	"unprotect_ack":1007,	// 
	"free_req":8,	// 结束攻击
	"free_ack":1008,	// 
	"getinfo_req":9,	// 读取玩家信息
	"getinfo_ack":1009,	// 
	"getofflinemsg_req":10,	// 获取离线信息
	"getofflinemsg_ack":1010,	// 
	"adduser_req":11,	// 注册一个新注册的玩家
	"adduser_ack":1011,	// 
}

var RCode map[int16]string = map[int16]string {
	0:"forward_req",
	1000:"forward_ack",
	1:"login_req",
	1001:"login_ack",
	2:"logout_req",
	1002:"logout_ack",
	3:"changescore_req",
	1003:"changescore_ack",
	4:"getlist_req",
	1004:"getlist_ack",
	5:"raid_req",
	1005:"raid_ack",
	6:"protect_req",
	1006:"protect_ack",
	7:"unprotect_req",
	1007:"unprotect_ack",
	8:"free_req",
	1008:"free_ack",
	9:"getinfo_req",
	1009:"getinfo_ack",
	10:"getofflinemsg_req",
	1010:"getofflinemsg_ack",
	11:"adduser_req",
	1011:"adduser_ack",
}

var ProtoHandler map[uint16]func(int32, *packet.Packet) ([]byte, error) = map[uint16]func(int32, *packet.Packet)([]byte, error){
	0:P_forward_req,
	1:P_login_req,
	2:P_logout_req,
	3:P_changescore_req,
	4:P_getlist_req,
	5:P_raid_req,
	6:P_protect_req,
	7:P_unprotect_req,
	8:P_free_req,
	9:P_getinfo_req,
	10:P_getofflinemsg_req,
	11:P_adduser_req,
}