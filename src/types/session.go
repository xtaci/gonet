package types

import "time"

import (
	"types/estate"
	"types/grid"
)

type Session struct {
	MQ          chan interface{} // Player's Internal Message Queue
	Basic       *Basic           //Basic Info
	Res         *Res             // Resource table
	Bitmap      *grid.Grid
	Estate      *estate.Manager
	Moves       []estate.Move
	HeartBeat   time.Time
	IsLoggedOut bool // represents user logout or connection failure
}
