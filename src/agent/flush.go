package main

import (
	"time"
)

import (
	"db/data_tbl"
	"db/user_tbl"
	. "types"
	"types/estates"
	"types/heroes"
	"types/samples"
	"types/soldiers"
)

//------------------------------------------------ flush all user data
func _flush(sess *Session) {
	user_tbl.Set(sess.User)
	data_tbl.Set(estates.COLLECTION, &sess.Estates)
	data_tbl.Set(soldiers.COLLECTION, &sess.Soldiers)
	data_tbl.Set(heroes.COLLECTION, &sess.Heroes)
	data_tbl.Set(samples.COLLECTION, &sess.LatencySamples)
	sess.LastFlushTime = time.Now().Unix()
	sess.OpCount = 0
}
