package core

import (
	"sync"
	"strconv"
	"time"
)

import (
	"misc/timer"
	"cfg"
)

type Record struct{
	_stats []*StatsObject
	_lock sync.Mutex
}

func (r *Record) Lock() {
	r._lock.Lock()
}

func (r *Record) Unlock() {
	r._lock.Unlock()
}

var (
	_stats map[int32]*Record
	_stats_lock sync.RWMutex
	_stats_chan chan int32
)

func init() {
	_stats = make(map[int32]*Record)
	_stats_chan = make(chan int32)
	go _writer()
}

func _writer() {
	for {
		user_id := <-_stats_chan
		collection := _stats[user_id]
	}
}

//------------------------------------------------ Group StatsObject by user 
func Collect(obj *StatsObject) {
	_stats_lock.RLock()
	record := _stats[obj.UserId]
	_stats_lock.RUnlock()

	config := cfg.Get()
	collect_interval, _ := strconv.Atoi(config["collect_interval"])

	if record != nil {
		record.Lock()
		if record._stats == nil {
			timer.Add(obj.UserId, time.Now().Unix()+int64(collect_interval), _stats_chan)
		}
		record._stats = append(record._stats, obj)
		record.Unlock()
	}
}
