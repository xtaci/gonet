package types

import (
	"encoding/json"
	"misc/crypto/pike"
	"net"
	"time"
)

import (
	"types/estates"
	"types/grid"
	"types/heroes"
	"types/samples"
	"types/soldiers"
)

type IPCObject struct {
	SrcID      int32  // sender id
	DestID     int32  // destination id
	Multicast  bool   // indicate wheather this message should be deliver to a group.
	Service    int16  // service type
	Object     []byte // json formatted object
	Time       int64  // sent time
	MarkDelete bool   // for db mark as delete
}

func (obj *IPCObject) Json() []byte {
	val, _ := json.Marshal(obj)
	return val
}

type Session struct {
	IP     net.IP
	MQ     chan IPCObject // Player's Internal Message Queue
	Crypto *pike.Pike     // a crypto algorithms
	// user data
	User     *User
	Estates  *estates.Manager
	Soldiers *soldiers.Manager
	Heroes   *heroes.Manager
	Grid     *grid.Grid // Building's bitmap, online constructing...
	Events   []int32    // event ids

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
