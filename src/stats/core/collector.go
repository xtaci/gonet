package core

import (
	"strconv"
	"sync"
	"time"
)

import (
	"cfg"
	. "db"
	"misc/timer"
)

type Record struct {
	_stats []*StatsObject
	_lock  sync.Mutex
}

func (r *Record) Lock() {
	r._lock.Lock()
}

func (r *Record) Unlock() {
	r._lock.Unlock()
}

//------------------------------------------------ variables
var (
	_stats      map[int32]*Record
	_stats_lock sync.RWMutex
	_stats_chan chan int32
)

const (
	MAX_WRITE_REQ    = 1000000
	STATS_COLLECTION = "STATS"
)

func init() {
	_stats = make(map[int32]*Record)
	_stats_chan = make(chan int32, MAX_WRITE_REQ)
	go _writer()
}

//------------------------------------------------ statistical data writer
func _writer() {
	for {
		user_id := <-_stats_chan
		record := _stats[user_id]
		if record != nil {
			_create_report(record)
		}
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

func _create_report(record *Record) {
	record.Lock()
	defer record.Unlock()
	// TODO: create a summary report
	config := cfg.Get()
	collection := Mongo.DB(config["mongo_db_stats"]).C(STATS_COLLECTION)
	// empty the stats
	record._stats = nil
}
