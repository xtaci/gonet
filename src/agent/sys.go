package main

import (
	"runtime"
	"time"
)

import (
	"agent/gsdb"
	. "helper"
	"misc/timer"
	. "types"
)

//---------------------------------------------------------- 系统routine,用户ID为0
func SysRoutine() {
	var sess Session
	sess.MQ = make(chan IPCObject, SYS_MQ_SIZE)
	gsdb.RegisterOnline(&sess, SYS_USR)

	// 定时器组
	gc_timer := make(chan int32, 10)
	// 初始触发
	gc_timer <- 1

	for {
		select {
		case msg, ok := <-sess.MQ: // IPC消息
			if !ok {
				return
			}
			// 只处理消息, 没有客户端可以返回
			IPCRequestProxy(&sess, &msg)
		case <-gc_timer: // 强制垃圾回收并打印性能日志
			runtime.GC()
			INFO("== PERFORMANCE LOG ==")
			INFO("Goroutine Count:", runtime.NumGoroutine())
			INFO("GC Summary:", GCSummary())
			INFO("Sysroutine MQ size:", len(sess.MQ))
			timer.Add(0, time.Now().Unix()+GC_INTERVAL, gc_timer)
		}
	}
}
