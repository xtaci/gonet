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
	"misc/timer"
)

//------------------------------------------------ variables
var (
	_all      map[int32]*Collector
	_all_lock sync.RWMutex
	CH        chan int32
)

const (
	DAY_SEC = int64(86400)
)

func init() {
	_all = make(map[int32]*Collector)
	CH = make(chan int32, 5)
	go _writer()

	config := cfg.Get()
	trigger, err := strconv.Atoi(config["collect_time"])
	if err != nil {
		log.Println("cannot parse collect_time from config", err)
	}

	// 寻找最近的触发点
	now := time.Now().Unix()
	passed := now % DAY_SEC
	if passed < int64(trigger) {
		timer.Add(-1, now-passed+int64(trigger), CH)
	} else {
		timer.Add(-1, now-passed+int64(trigger)+DAY_SEC, CH)
	}
}

//------------------------------------------------ 统计数据定时汇总写入
func _writer() {
	for {
		// 时钟信号
		<-CH

		// 复制map已进行费事操作,不阻塞collect
		_all_lock.Lock()
		snapshot := make(map[int32]*Collector)
		for k, v := range _all {
			snapshot[k] = v
		}
		_all_lock.Unlock()

		// 创建每个玩家的报表
		c := StatsCollection()
		for userid, collector := range snapshot {
			if collector != nil {
				archive := _archive(userid, collector)
				c.Upsert(bson.M{"userid": userid}, archive)
			}
		}

		now := time.Now().Unix()
		log.Printf("stats flush finished at %v\n", now)

		// 明天同一时刻再见
		passed := now % DAY_SEC
		config := cfg.Get()
		trigger, err := strconv.Atoi(config["collect_time"])
		if err != nil {
			log.Println("cannot parse collect_time from config", err)
		}

		timer.Add(-1, now-passed+int64(trigger)+DAY_SEC, CH)
		snapshot = nil
	}
}

//------------------------------------------------ 按用户分组消息收集
func Collect(userid int32, obj *StatsObject) {
	// 获得记录,为空则创建
	_all_lock.Lock()
	collector := _all[userid]
	if collector == nil {
		collector = &Collector{}
		_all[userid] = collector
	}
	_all_lock.Unlock()
	_drop_expired(collector)

	// 放入新的消息
	collector.Lock()
	defer collector.Unlock()
	if collector._stats == nil {
		collector._stats = make([]*StatsObject, 0, 512)
	}
	collector._stats = append(collector._stats, obj)
}
