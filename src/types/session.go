package types

import "time"

type Session struct {
	MQ chan interface{}
	ServerMQ chan []byte

	User   User

	SESSID    [128]byte // UNIQUE session ID
	HeartBeat time.Time
}
