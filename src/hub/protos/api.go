package protos

import "misc/packet"

var Code map[string]int16 = map[string]int16 {
	"forward_req":0,	// payload:msg 消息转发
	"login_req":1,	// payload:id 登陆
	"login_ack":1001,	// payload:command_result_pack 
	"logout_req":2,	// payload:id 登出
	"logout_ack":1002,	// payload:command_result_pack 登出
	"changescore_req":3,	// payload:changescore 改变分数
	"changescore_ack":1003,	// payload:command_result_pack 改变分数
	"getlist_req":4,	// payload:getlist 获取列表
	"getlist_ack":1004,	// payload:getlist_result 获取列表
	"raid_req":5,	// payload:id 攻击
	"raid_ack":1005,	// payload:command_result_pack 攻击
	"protect_req":6,	// payload:id 加保护
	"protect_ack":1006,	// payload:command_result_pack 加保护
	"unprotect_req":7,	// payload:id 撤销保护
	"unprotect_ack":1007,	// payload:command_result_pack 撤销保护
	"free_req":8,	// payload:id 结束攻击
	"free_ack":1008,	// payload:command_result_pack 结束攻击
	"getstate_req":9,	// payload:id 读取状态
	"getstate_ack":1009,	// payload:command_result_pack 读取状态
	"getprotecttime_req":10,	// payload:id 获取保护时间截止
	"getprotecttime_ack":1010,	// payload:timeresult 获取保护时间截止
	"getname_req":11,	// payload:id 获取玩家名字
	"getname_ack":1011,	// payload:stringresult 获取玩家名字
}

var RCode map[int16]string = map[int16]string {
	0:"forward_req",
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
	9:"getstate_req",
	1009:"getstate_ack",
	10:"getprotecttime_req",
	1010:"getprotecttime_ack",
	11:"getname_req",
	1011:"getname_ack",
}

var ProtoHandler map[uint16]func(int32, *packet.Packet) ([]byte, error) = map[uint16]func(int32, *packet.Packet)([]byte, error){
	0:_forward_req,
	1:_login_req,
	2:_logout_req,
	3:_changescore_req,
	4:_getlist_req,
	5:_raid_req,
	6:_protect_req,
	7:_unprotect_req,
	8:_free_req,
	9:_getstate_req,
	10:_getprotecttime_req,
	11:_getname_req,
}