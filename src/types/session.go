package types

import "time"

import (
	"types/grid"
)

type Session struct {
	MQ chan interface{} // Player's Internal Message Queue

	User        User
	Bitmap      *grid.Grid
	Data		*PlayerData
	HeartBeat   time.Time
	IsLoggedOut bool // represents user logout or connection failure
}
