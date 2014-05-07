package main

import (
	"log"
	"sync"
	"time"
)

import (
	"helper"
	"misc/timer"
	. "types"
)

func init() {
	log.SetPrefix("[GS] ")
}

var wg sync.WaitGroup
var die chan bool // for server close

func init() {
	die = make(chan bool)
}

//----------------------------------------------- Start Agent when a client is connected
func StartAgent(sess *Session, in chan []byte, out *Buffer) {
	wg.Add(1)
	defer wg.Done()
	defer helper.PrintPanicStack()

	// init session
	sess.MQ = make(chan IPCObject, DEFAULT_MQ_SIZE)
	sess.ConnectTime = time.Now()
	sess.LastPacketTime = time.Now()

	// custom-sec timer, 60-sec
	custom_timer := make(chan int32, 1)
	timer.Add(-1, time.Now().Unix()+CUSTOM_TIMER, custom_timer)

	// cleanup work
	defer func() {
		close_work(sess)
	}()

	// the main message loop
	for {
		select {
		case msg, ok := <-in: // network protocol
			if !ok {
				return
			}

			sess.PacketTime = time.Now()
			if result := UserRequestProxy(sess, msg); result != nil {
				err := out.Send(result)
				if err != nil {
					helper.ERR("cannot send to client", err)
					return
				}
			}
			sess.LastPacketTime = sess.PacketTime
			sess.PacketCount++ // packet count
		case msg := <-sess.MQ: // internal IPC
			if result := IPCRequestProxy(sess, &msg); result != nil {
				err := out.Send(result)
				if err != nil {
					helper.ERR("cannot send ipc response", err)
					return
				}
			}
		case <-custom_timer: // 60-sec timer
			timer_work(sess)
			timer.Add(-1, time.Now().Unix()+CUSTOM_TIMER, custom_timer)
		case <-die:
			sess.Flag |= SESS_KICKED_OUT
		}

		// is the session been kicked out
		if sess.Flag&SESS_KICKED_OUT != 0 {
			return
		}
	}
}
