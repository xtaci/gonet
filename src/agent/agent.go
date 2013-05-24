package agent

import (
	"fmt"
	"net"
	"time"
)

import (
	. "types"
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
	sess.MQ = make(chan interface{}, 128)

	// session timeout
	session_timeout := make(chan bool)
	go _timer(5, session_timeout)

	// event_timeout(such as, upgrades...)
	event_timer := make(chan bool)
	go _timer(1, event_timer)

	// write buffer
	bufctrl := make(chan bool)
	buf := NewBuffer(conn, bufctrl)
	go buf.Start()

	// cleanup work
	defer func() {
		close_work(&sess)
		close(session_timeout)
		close(event_timer)
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

			sess.HeartBeat = time.Now().Unix()
			if result := UserRequestProxy(&sess, msg); result != nil {
				fmt.Println(result)
				err := buf.Send(result)
				if err != nil {
					return
				}
			}
		case msg, ok := <-sess.MQ: // async
			if !ok {
				return
			}

			if result := IPCRequestProxy(&sess, msg); result != nil {
				fmt.Println(result)
				err := buf.Send(result)
				if err != nil {
					return
				}
			}

		case _ = <-session_timeout:
			if session_work(&sess) {
				return
			}
		case _ = <-event_timer:
			event_work(&sess)
		}
	}
}
