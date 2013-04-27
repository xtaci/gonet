package types

import "time"

const (
	FREE = iota
	ONLINE
	RAID // being raid
)

type Session struct {
	MQ     chan interface{}
	User   User
	Cities []City

	SESSID    [128]byte // UNIQUE session ID
	HeartBeat time.Time

	Status	int
}
