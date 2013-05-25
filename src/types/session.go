package types

import (
	"types/estate"
	"types/grid"
)

type IPCObject struct {
	Sender  int32 // sender id
	Service int16
	Object  []byte //json formatted object
}

type Session struct {
	MQ            chan IPCObject // Player's Internal Message Queue
	Basic         *Basic         //Basic Info
	Res           *Res           // Resource table
	Bitmap        *grid.Grid
	EstateManager *estate.Manager
	Moves         []estate.Move
	HeartBeat     int64
	IsLoggedOut   bool // represents user logout or connection failure
}
