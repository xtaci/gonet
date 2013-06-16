package core

import (
	"encoding/json"
	"sync"
)

import (
	. "types"
	"types/estates"
)

const (
	UNKNOWN  = int32(iota)
	TYPE_SUM // sum the result
)

type StatsObject struct {
	Type      int32
	Key       string
	Value     float32
	Timestamp int64
}

//------------------------------------------------ 一个玩家对应一个
type Collector struct {
	_stats []*StatsObject
	_lock  sync.Mutex
}

func (r *Collector) Lock() {
	r._lock.Lock()
}

func (r *Collector) Unlock() {
	r._lock.Unlock()
}

type Archive struct {
	UserId    int32
	Timestamp int64
	Fields    map[string]float32
	User      *User
	Estates   estates.Manager
}

func (archive *Archive) Marshal() []byte {
	json_val, _ := json.Marshal(archive)
	return json_val
}
