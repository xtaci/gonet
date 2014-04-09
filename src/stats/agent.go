package main

import (
	"log"
	"net"
	"time"
)

import (
	"misc/packet"
	"misc/timer"
)

const (
	DEFAULT_MAX_QUEUE = 1024 * 1000
	PRINT_INTERVAL    = 300
)

func init() {
	log.SetPrefix("[STATS] ")
}

//------------------------------------------------ Stats Server Agent
func StatsAgent(incoming chan []byte, conn net.Conn) {
	queue_timer := make(chan int32, 1)
	queue_timer <- 1

	for {
		select {
		case sample := <-incoming:
			reader := packet.Reader(sample)
			HandleRequest(reader)
		case <-queue_timer:
			log.Println("============== STATS QUEUE SIZE:", len(incoming), "===================")
			timer.Add(1, time.Now().Unix()+PRINT_INTERVAL, queue_timer)
		}
	}
}
