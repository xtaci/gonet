package types

import (
	"types/estates"
	"types/grid"
)

type IPCObject struct {
	Sender  int32 // sender id
	Service int16
	Object  []byte // json formatted object
	Time    int64  // send time
}

type Session struct {
	MQ chan IPCObject // Player's Internal Message Queue
	// user data
	User    User
	Estates estates.Manager
	Grid    *grid.Grid // Building's bitmap, online constructing...

	// session related
	LoggedIn bool // flag for weather the user is logged in
	KickOut  bool // flag for player is kicked out

	// time related
	ConnectTime    int64 // tcp connection establish time
	LastPacketTime int64 // last packet arrive time
	LastFlushTime  int64 // last flush to db time
	OpCount        int   // num of operations since last sync
}
