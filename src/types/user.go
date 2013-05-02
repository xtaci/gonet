package types

import (
	"time"
	"sync/atomic"
)

const (
	FREE = iota
	ONLINE
	BEING_RAID
	PROTECTED
)

type User struct {
	Id           int32
	Name         string
	Mac          string
	Score        int32
	State        int32
	Archives     string
	LastSaveTime time.Time
	ProtectTime  time.Time
	CreatedAt    time.Time
}

func (ud *User) ChangeState(oldstate, newstate int32) bool {
	return atomic.CompareAndSwapInt32(&ud.State, oldstate, newstate)
}
