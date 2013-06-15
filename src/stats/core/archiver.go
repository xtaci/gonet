package core

import (
	"db/data_tbl"
	"types/estates"
)

//------------------------------------------------ 归档玩家同级数据
func _archive(userid int32, record *Record) *Archive {
	_drop_expired(record)
	// TODO: create a summary report within last 24-hours
	record.Lock()
	defer record.Unlock()

	archive := &Archive{}
	archive.Fields = make(map[string]string)
	manager := &estates.Manager{}
	data_tbl.Get(estates.COLLECTION, userid, manager)

	return archive
}
