package types

import (
	"types/estate"
	"types/grid"
)

type IPCObject struct {
	Sender  int32 // sender id
	Service int16
	Object  []byte // json formatted object
	Time    int64  // send time
}

type Session struct {
	MQ            chan IPCObject // Player's Internal Message Queue
	Basic         Basic          //Basic Info
	Res           Res            // Resource table
	Bitmap        grid.Grid
	EstateManager estate.Manager
	Moves         []estate.Move
	IsLoggedOut   bool // represents user logout or connection failure

	ConnectTime int64
	LastPing    int64
	LastFlush   int64 // last flush to db time
	OpCount     int   // num of operations since last sync
}
