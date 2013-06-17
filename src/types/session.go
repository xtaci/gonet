package types

import (
	"encoding/json"
	"net"
	"time"
)

import (
	"types/estates"
	"types/grid"
	"types/samples"
	"types/soldiers"
)

type IPCObject struct {
	Sender    int32 // sender id
	Multicast bool  // indicate wheather this message is a multicast message(group)
	Service   int16
	Object    []byte // json formatted object
	Time      int64  // send time
}

func (obj *IPCObject) Json() []byte {
	val, _ := json.Marshal(obj)
	return val
}

type Session struct {
	IP net.IP
	MQ chan IPCObject // Player's Internal Message Queue
	// user data
	User     *User
	Estates  *estates.Manager
	Soldiers *soldiers.Manager
	Grid     *grid.Grid // Building's bitmap, online constructing...

	// session related
	LoggedIn bool // flag for weather the user is logged in
	KickOut  bool // flag for player is kicked out

	// time related
	ConnectTime    time.Time        // tcp connection establish time, in millsecond(ms)
	LastPacketTime int64            // last packet arrive time, in seconds(s)
	LastFlushTime  int64            // last flush to db time, in seconds(s)
	OpCount        int              // num of operations since last sync
	LatencySamples *samples.Samples // 网络延迟样本
}
