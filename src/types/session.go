package types

import "time"

import (
	"types/army"
	"types/defense"
	"types/grid"
)

type Session struct {
	MQ   chan interface{} // Player's Internal Message Queue

	User      User
	Bitmap    *grid.Grid
	Army	  []army.Building
	Defense   []defense.Building
	HeartBeat time.Time
	IsLoggedOut bool	// represents user logout or connection failure
}
