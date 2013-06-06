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
	. "db"
	"misc/timer"
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
	DAY_SEC          = int64(86400)
	STATS_COLLECTION = "STATS"
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
		c := Mongo.DB(config["mongo_db"]).C(STATS_COLLECTION)

		// 复制map已进行费事操作,不阻塞collect
		_all_lock.Lock()
		snapshot := make(map[int32]*Record)
		for k, v := range _all {
			snapshot[k] = v
		}
		_all_lock.Unlock()

		for _, record := range snapshot {
			if record != nil {
				summary := _create_summary(record)
				c.Upsert(bson.M{"userid": summary.UserId}, summary)
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
func Collect(obj *StatsObject) {
	// 获得记录,为空则创建
	_all_lock.Lock()
	record := _all[obj.UserId]
	if record == nil {
		record = &Record{}
		_all[obj.UserId] = record
	}
	_all_lock.Unlock()

	// 丢弃过期消息
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
	record._stats = append(record._stats, obj)
	record.Unlock()
}

//------------------------------------------------ create a summary and remove old data
func _create_summary(record *Record) *Summary {
	record.Lock()
	defer record.Unlock()
	// TODO: create a summary report
	sum := &Summary{}
	// empty the stats
	record._stats = nil

	return sum
}
