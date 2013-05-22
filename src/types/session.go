package types

import "time"

import (
	"types/grid"
	"types/estate"
)

type Session struct {
	MQ chan interface{} // Player's Internal Message Queue
	Basic		Basic //Basic Info
	Bitmap      *grid.Grid
	Estate		*estate.EstateManager
	HeartBeat   time.Time
	IsLoggedOut bool // represents user logout or connection failure
}
