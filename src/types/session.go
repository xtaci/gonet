package types

import (
	"types/basic"
	"types/defensive"
	"types/grid"
	"types/offensive"
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
	Res Res
	// Building Info
	BAS    basic.Manager
	DEF    defensive.Manager
	OFS    offensive.Manager
	Bitmap grid.Grid // Building's bitmap

	// Session Info
	LoggedIn bool // flag for weather the user is logged in

	ConnectTime int64
	LastPing    int64
	LastFlush   int64 // last flush to db time
	OpCount     int   // num of operations since last sync
	KickOut     bool  // flag for is kicked out
}
