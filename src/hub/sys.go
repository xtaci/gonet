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
	gc_timer := make(chan int32, 10)
	gc_timer <- 1

	for {
		select {
		case <-gc_timer:
			// gc work
			runtime.GC()
			log.Println("GC executed")
			log.Println("NumGoroutine", runtime.NumGoroutine())
			log.Println("GC Summary:", helper.GCSummary())
			timer.Add(0, time.Now().Unix()+int64(GC_INTERVAL), gc_timer)
		}
	}
}
