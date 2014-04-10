package types

import (
	"fmt"
	"misc/crypto/pike"
	"net"
	"time"
)

import (
// TODO: import data structure
)

const (
	SESS_LOGGED_IN  = 0x1
	SESS_KICKED_OUT = 0x2
	SESS_REGISTERED = 0x4
	SESS_KEYEXCG    = 0x8
	SESS_ENCRYPT    = 0x10
)

type Session struct {
	IP      net.IP
	MQ      chan IPCObject // Player's Internal Message Queue
	Encoder *pike.Pike
	Decoder *pike.Pike

	// TODO: all user data structure
	User *User

	// session related
	LoggedIn bool // flag for weather the user is logged in
	KickOut  bool // flag for player is kicked out

	// session flag
	Flag int32

	// time related variables
	ConnectTime    time.Time // tcp connection establish time, in millsecond(ms)
	PacketTime     time.Time // last packet time
	LastPacketTime time.Time // last packet arrive time, in seconds(s)
	_dirtycount    int32     // dirty ops count

	// packet rate control
	PacketCount int64 // count packets
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
