package agent

import (
	"log"
)

import (
	"db/city"
	"db/user"
	. "types"
)

func db_work(sess *Session) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic when flushing database: %v", x)
		}
	}()

	if sess.User.Id != 0 {
		user.Flush(&sess.User)
		for i := range sess.Cities {
			city.Flush(&sess.Cities[i])
		}
	}
}

