package types

import "time"

type Session struct {
	MQ     chan interface{}
	User   User
	Cities []City

	SESSID    [128]byte // UNIQUE session ID
	HeartBeat time.Time
}
