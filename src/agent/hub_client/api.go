package hub_client

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
