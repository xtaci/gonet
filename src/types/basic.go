package types

import (
	"encoding/json"
	"sync/atomic"
)

type Basic struct {
	Id	int32
	Name string
	Score int32
	ProtectTimeout int64
	IsProtecting bool
	LoginCount	int32
	LastLogin	int64
	NextVal	int32
}

func (b *Basic) JSON() string {
	val, _ := json.Marshal(b)
	return string(val)
}

func (b *Basic) GENID() int32 {
	return atomic.AddInt32(&b.NextVal, 1)
}
