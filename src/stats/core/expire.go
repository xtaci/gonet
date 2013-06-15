package core

import (
	"time"
)

//------------------------------------------------ 丢弃过期统计数据
func _drop_expired(collector *Collector) {
	collector.Lock()
	defer collector.Unlock()

	expire_point := time.Now().Unix() - DAY_SEC
	count := 0
	for _, v := range collector._stats {
		if v.Timestamp < expire_point {
			count++
		} else {
			break
		}
	}

	if count > 0 {
		collector._stats = collector._stats[count:]
	}
}
