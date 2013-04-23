package types

import "time"

type Session struct {
	MQ     chan string
	User   User
	Cities []City

	HeartBeat time.Time
}
