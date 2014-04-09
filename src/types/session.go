package types

import (
	"fmt"
	"misc/crypto/pike"
	"net"
	"time"
)

import (
	"types/estates"
	"types/grid"
	"types/heroes"
	"types/soldiers"
)

const (
	SESS_LOGGED_IN  = 0x1  // 玩家是否登录
	SESS_KICKED_OUT = 0x2  // 玩家是否被服务器踢掉
	SESS_REGISTERED = 0x4  // 是否已经执行过注册(避免恶意注册)
	SESS_KEYEXCG    = 0x8  // 是否已经交换完毕KEY
	SESS_ENCRYPT    = 0x10 // 是否可以开始加密
)

type Session struct {
	IP      net.IP
	MQ      chan IPCObject // Player's Internal Message Queue
	Encoder *pike.Pike     // 加密器
	Decoder *pike.Pike     // 解密器

	// user data
	User     *User
	Estates  *estates.Manager
	Soldiers *soldiers.Manager
	Heroes   *heroes.Manager
	Grid     *grid.Grid // Building's bitmap, online constructing...

	// session related
	LoggedIn bool // flag for weather the user is logged in
	KickOut  bool // flag for player is kicked out

	// 会话标记
	Flag int32
	// time related
	ConnectTime    time.Time // tcp connection establish time, in millsecond(ms)
	PacketTime     time.Time // 当前包的到达时间
	LastPacketTime time.Time // last packet arrive time, in seconds(s)
	LastFlushTime  time.Time // last flush to db time, in seconds(s)
	_dirtycount    int32     // 脏数据计数
	// RPS控制
	PacketCount int64 // 对收到的包进行计数，避免恶意发包
}

func (sess *Session) MarkDirty() {
	sess._dirtycount++
}

func (sess *Session) DirtyCount() int32 {
	return sess._dirtycount
}

func (sess *Session) MarkClean() {
	sess._dirtycount = 0
}

//------------------------------------------------ integer to string
func S(a interface{}) string {
	return fmt.Sprint(a)
}
