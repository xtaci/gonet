package types

import "time"

type Session struct {
	MQ   chan interface{} // ASYNC
	CALL chan interface{} // SYNC

	User      User
	Bitmap    []byte
	Buildings []byte
	HeartBeat time.Time
}
