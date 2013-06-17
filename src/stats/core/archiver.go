package core

import (
	"log"
)

import (
	"db/data_tbl"
	"db/user_tbl"
	"types/estates"
)

//------------------------------------------------ 归档玩家数据
func _archive(userid int32, collector *Collector) *Archive {
	_drop_expired(collector)

	collector.Lock()
	defer collector.Unlock()

	archive := &Archive{}
	archive.Fields = make(map[string]float32)

	for _, stat := range collector._stats {
		switch stat.Type {
		case TYPE_SUM:
			archive.Fields[stat.Key] += stat.Value
		default:
			log.Println("cannot archive:", stat)
		}
	}

	// snapshot of player data
	data_tbl.Get(estates.COLLECTION, userid, &archive.Estates)
	archive.User = user_tbl.Get(userid)
	return archive
}
