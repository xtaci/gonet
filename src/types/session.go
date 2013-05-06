package types

import "time"

type Session struct {
	MQ      chan interface{}		// ASYNC
	CALL	chan interface{}		// SYNC

	User User

	SESSID    [128]byte // UNIQUE session ID
	HeartBeat time.Time
}
