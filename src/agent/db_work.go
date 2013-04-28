package agent

import (
	"log"
)

import (
	"db/city_tbl"
	"db/user_tbl"
	. "types"
)

func db_work(sess *Session) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic when flushing database: %v", x)
		}
	}()

	if sess.User.Id != 0 {
		user_tbl.Flush(&sess.User)
		for i := range sess.Cities {
			city_tbl.Flush(&sess.Cities[i])
		}
	}
}

