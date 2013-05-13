package types

import "time"

type Session struct {
	MQ   chan interface{} // Player's Internal Message Queue

	User      User
	Bitmap    []byte
	Buildings []byte
	HeartBeat time.Time
}
