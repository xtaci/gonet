package main

import (
	"log"
	"runtime"
	"time"
)

import (
	"agent/gsdb"
	"cfg"
	"misc/timer"
	"strconv"
	. "types"
)

const (
	SYS_USR     = 0
	SYS_MQ_SIZE = 65535
	GC_INTERVAL = 300
)

//---------------------------------------------------------- 系统routine,用户ID为0
func SysRoutine() {
	var sess Session
	sess.MQ = make(chan IPCObject, SYS_MQ_SIZE)
	gsdb.RegisterOnline(&sess, SYS_USR)

	// timer
	gc_timer := make(chan int32, 1)
	timer.Add(0, time.Now().Unix()+GC_INTERVAL, gc_timer)

	for {
		config := cfg.Get()
		gc_interval, e := strconv.Atoi(config["gc_interval"])
		if e != nil {
			gc_interval = GC_INTERVAL
		}

		select {
		case msg, ok := <-sess.MQ: // async
			if !ok {
				return
			}
			IPCRequestProxy(&sess, &msg)
		case <-gc_timer:
			runtime.GC()
			log.Println("GC executed")
			log.Println("NumGoroutine", runtime.NumGoroutine())
			timer.Add(0, time.Now().Unix()+int64(gc_interval), gc_timer)
		}
	}
}
