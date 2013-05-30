package agent

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

import (
	"cfg"
	. "types"
)

const (
	DEFAULT_MQ_SIZE = 128
)

//----------------------------------------------- a simple timer
func _timer(interval int, ch chan bool) {
	defer func() {
		recover()
	}()

	for {
		time.Sleep(time.Duration(interval) * time.Second)
		ch <- true
	}
}

//----------------------------------------------- Start Agent when a client is connected
func StartAgent(in chan []byte, conn net.Conn) {
	var sess Session
	sess.MQ = make(chan IPCObject, DEFAULT_MQ_SIZE)
	sess.ConnectTime = time.Now().Unix()
	sess.LastPacketTime = time.Now().Unix()

	// session timeout
	session_timeout := make(chan bool)
	go _timer(5, session_timeout)

	// event_timeout(such as, upgrade, flush)
	std_timer := make(chan bool)
	go _timer(1, std_timer)

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
		close(session_timeout)
		close(std_timer)
		close(sess.MQ)
		bufctrl <- false
		conn.Close()
	}()

	// the main message loop
	for {
		select {
		case msg, ok := <-in:
			if !ok {
				return
			}

			if result := UserRequestProxy(&sess, msg); result != nil {
				fmt.Println(result)
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
		case _ = <-session_timeout:
			if session_work(&sess) {
				return
			}
		case _ = <-std_timer:
			timer_work(&sess)
		}

		// TODO: 持久化逻辑#1： 超过一定的操作数量，刷入数据库
		if sess.OpCount > flush_ops {
			sess.OpCount = 0
		}
	}
}
