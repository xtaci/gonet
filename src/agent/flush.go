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
	if sess.Estates != nil {
		data_tbl.Set(estates.COLLECTION, sess.Estates)
	}

	if sess.Soldiers != nil {
		data_tbl.Set(soldiers.COLLECTION, sess.Soldiers)
	}

	if sess.Heroes != nil {
		data_tbl.Set(heroes.COLLECTION, sess.Heroes)
	}

	if sess.LatencySamples != nil {
		data_tbl.Set(samples.COLLECTION, sess.LatencySamples)
	}

	sess.LastFlushTime = time.Now().Unix()
	sess.OpCount = 0
}
