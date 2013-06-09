package main

import (
	"log"
	"net"
	"strconv"
	"time"
)

import (
	"cfg"
	"misc/timer"
	. "types"
)

const (
	DEFAULT_MQ_SIZE = 128
)

func init() {
	log.SetPrefix("[GS]")
}

//----------------------------------------------- Start Agent when a client is connected
func StartAgent(in chan []byte, conn net.Conn) {
	var sess Session
	sess.MQ = make(chan IPCObject, DEFAULT_MQ_SIZE)
	sess.ConnectTime = time.Now()
	sess.LastPacketTime = time.Now().Unix()
	sess.KickOut = false

	// standard 1-sec timer
	std_timer := make(chan int32, 1)
	timer.Add(1, time.Now().Unix()+1, std_timer)

	// write buffer
	bufctrl := make(chan bool)
	buf := NewBuffer(conn, bufctrl)
	go buf.Start()

	// max #operartion before flush
	config := cfg.Get()
	flush_ops, _ := strconv.Atoi(config["flush_ops"])

	// cleanup work
	defer func() {
		close_work(&sess)
		close(std_timer)
		close(sess.MQ)
		bufctrl <- false
	}()

	// the main message loop
	for {
		select {
		case msg, ok := <-in:
			if !ok {
				return
			}

			if result := UserRequestProxy(&sess, msg); result != nil {
				err := buf.Send(result)
				if err != nil {
					return
				}
			}
			sess.LastPacketTime = time.Now().Unix()

		case msg, ok := <-sess.MQ: // async
			if !ok {
				return
			}
			IPCRequestProxy(&sess, &msg)

		case _ = <-std_timer:
			timer_work(&sess)
			if session_timeout(&sess) {
				return
			}
			timer.Add(1, time.Now().Unix()+1, std_timer)
		}

		// TODO: 持久化逻辑#1： 超过一定的操作数量，刷入数据库
		if sess.OpCount > flush_ops {
			sess.OpCount = 0
			sess.Dirty = false
		}

		// 是否被逻辑踢出
		if sess.KickOut {
			return
		}
	}
}
