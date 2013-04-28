package types

import "time"

const (
	FREE = iota
	ONLINE
	BEING_RAID
)

type User struct {
	Id          int32
	Name        string
	Mac			string
	Score		int32
	State		int
	Archives	string
	LastSaveTime	time.Time
	ProtectTime	time.Time
	CreatedAt	time.Time
}
