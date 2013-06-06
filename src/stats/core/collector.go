package core

import (
	"labix.org/v2/mgo/bson"
	"log"
	"strconv"
	"sync"
	"time"
)

import (
	"cfg"
	"db/data_tbl"
	"misc/timer"
	"types/estates"
)

//------------------------------------------------ 一个玩家对应一个
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
	_all      map[int32]*Record
	_all_lock sync.RWMutex
	CH        chan int32
)

const (
	DAY_SEC = int64(86400)
)

func init() {
	_all = make(map[int32]*Record)
	CH = make(chan int32, 5)
	go _writer()

	config := cfg.Get()
	trigger, _ := strconv.Atoi(config["collect_time"])

	// 寻找最近的触发点
	now := time.Now().Unix()
	passed := now % DAY_SEC
	if passed < int64(trigger) {
		timer.Add(-1, now-passed+int64(trigger), CH)
	} else {
		timer.Add(-1, now-passed+int64(trigger)+DAY_SEC, CH)
	}
}

//------------------------------------------------ statistical data writer
func _writer() {
	for {
		// 时钟信号
		<-CH

		config := cfg.Get()

		// 复制map已进行费事操作,不阻塞collect
		_all_lock.Lock()
		snapshot := make(map[int32]*Record)
		for k, v := range _all {
			snapshot[k] = v
		}
		_all_lock.Unlock()

		c := StatsCollection()
		for userid, record := range snapshot {
			if record != nil {
				summary := _create_summary(userid, record)
				c.Upsert(bson.M{"userid": userid}, summary)
			}
		}

		now := time.Now().Unix()
		passed := now % DAY_SEC

		log.Printf("stats flush finished at %v\n", now)

		// 明天同一时刻再见
		trigger, _ := strconv.Atoi(config["collect_time"])
		timer.Add(-1, now-passed+int64(trigger)+DAY_SEC, CH)
		snapshot = nil
	}
}

//------------------------------------------------ 按用户分组消息收集
func Collect(userid int32, obj *StatsObject) {
	// 获得记录,为空则创建
	_all_lock.Lock()
	record := _all[userid]
	if record == nil {
		record = &Record{}
		_all[userid] = record
	}
	_all_lock.Unlock()
	_drop_expired(record)

	record.Lock()
	record._stats = append(record._stats, obj)
	record.Unlock()
}

//------------------------------------------------ create a summary and remove old data
func _create_summary(userid int32, record *Record) *Summary {
	record.Lock()
	defer record.Unlock()

	// TODO: create a summary report within last 24-hours
	sum := &Summary{}
	manager := &estates.Manager{}
	data_tbl.Get(estates.COLLECTION, userid, manager)

	_drop_expired(record)
	return sum
}

//------------------------------------------------ 丢弃过期消息
func _drop_expired(record *Record) {
	now := time.Now().Unix()
	record.Lock()
	count := 0
	for _, v := range record._stats {
		if v.Timestamp < now {
			count++
		} else {
			break
		}
	}

	if count > 0 {
		record._stats = record._stats[count:]
	}

	// 放入新的消息
	if record._stats == nil {
		record._stats = make([]*StatsObject, 0, 512)
	}
	record.Unlock()
}
