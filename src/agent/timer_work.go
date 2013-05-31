package agent

import (
	"fmt"
	"strconv"
	"time"
)

import (
	"cfg"
	"db/data_tbl"
	"db/user_tbl"
	. "types"
	"types/estates"
)

//----------------------------------------------- timer work
func timer_work(sess *Session) {
	// check whether the user is logged in
	if !sess.LoggedIn {
		return
	}

	// TODO: 持久化逻辑#2： 超过一定的时间，刷入数据库
	config := cfg.Get()
	ivl, _ := strconv.Atoi(config["flush_interval"])
	if time.Now().Unix()-sess.LastFlushTime > int64(ivl) {
		fmt.Println("TODO: flush all to db")
		flag1 := user_tbl.Set(&sess.User)
		flag2 := data_tbl.Set(estates.COLLECTION, &sess.Estates)

		// 成功后才设置
		if flag1 && flag2 {
			sess.LastFlushTime = time.Now().Unix()
		}
	}
}
