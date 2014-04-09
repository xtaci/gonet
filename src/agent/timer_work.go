package main

import (
	"strconv"
	"sync/atomic"
	"time"
)

import (
	"cfg"
	"helper"
	. "types"
)

//----------------------------------------------- 1分钟定时器，主要用于需要定时器的活动
func timer_work(sess *Session) {
	if sess.Flag&SESS_LOGGED_IN == 0 {
		return
	}

	// SIGTERM 信号检测
	if atomic.LoadInt32(&SIGTERM) == 1 {
		sess.Flag |= SESS_KICKED_OUT
		helper.NOTICE("收到SIGTERM, 玩家被动退出", sess.User.Id, sess.User.Name)
	}

	// 发包频率控制，太高的RPS直接踢掉
	config := cfg.Get()
	rpm_limit, _ := strconv.ParseFloat(config["rpm_limit"], 32)
	rpm := float64(sess.PacketCount) / float64(time.Now().Unix()-sess.ConnectTime.Unix()) * 60

	if rpm > rpm_limit {
		sess.Flag |= SESS_KICKED_OUT
		helper.ERR("玩家RPM太高", sess.User.Id, sess.User.Name, "RPM:", rpm)
		return
	}

	// 尝试持久化
	_flush_work(sess)
}
