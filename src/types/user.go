package types

import "time"

const (
	FREE = iota
	ONLINE
	RAID // being raid
)

type User struct {
	Id          int
	Name        string
	Mac			string
	Status      int
	Score		int
	LastSync	time.Time
	ShieldUntil time.Time
	CreatedAt   time.Time
}
