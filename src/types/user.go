package types

import (
	"sync/atomic"
)

type User struct {
	Id             int32
	ClanId         int32
	Name           string
	Pass           []byte
	Mac            string
	Score          int32
	ProtectTimeout int64
	LoginCount     int32
	LastLogin      int64
	NextVal        uint32
}

//------------------------------------------------ sequence generator
func (u *User) Next() uint32 {
	return atomic.AddUint32(&u.NextVal, 1)
}
