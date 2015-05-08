package main

import (
	"log"
	"runtime"
	"time"
)

import (
	"helper"
	"misc/timer"
)

const (
	GC_INTERVAL = 300
)

//---------------------------------------------------------- 系统routine
func SysRoutine() {
	// timer
	gc_timer := make(chan int32, 1)
	timer.Add(0, time.Now().Unix()+GC_INTERVAL, gc_timer)

	for {
		select {
		case <-gc_timer:
			runtime.GC()
			log.Println("GC executed")
			log.Println("NumGoroutine", runtime.NumGoroutine())
			log.Println("GC Summary:", helper.GCSummary())
			timer.Add(0, time.Now().Unix()+GC_INTERVAL, gc_timer)
		}
	}
}
