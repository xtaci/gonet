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
	MQ   chan IPCObject // Player's Internal Message Queue
	User User
	// Resources Info
	Res estates.ResManager
	// Building Info
	BAS    estates.BasManager
	DEF    estates.DefManager
	OFF		estates.OffManager
	Bitmap grid.Grid // Building's bitmap, online constructing...

	// Session Info
	LoggedIn bool // flag for weather the user is logged in

	ConnectTime int64
	LastPing    int64
	LastFlush   int64 // last flush to db time
	OpCount     int   // num of operations since last sync
	KickOut     bool  // flag for is kicked out
}
