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

//----------------------------------------------- timer work
func _timer(interval int, ch chan string) {
	defer func() {
		recover()
	}()

	for {
		time.Sleep(time.Duration(interval) * time.Second)
		ch <- "ding!dong!"
	}
}

//----------------------------------------------- Start Agent when a client is connected
func StartAgent(in chan []byte, conn net.Conn) {
	var sess Session
	sess.MQ = make(chan interface{}, 128)

	config := cfg.Get()

	// session timeout
	timer_ch_session := make(chan string)
	session_timeout := 30 // sec
	if config["session_timeout"] != "" {
		session_timeout, _ = strconv.Atoi(config["session_timeout"])
	}

	go _timer(session_timeout, timer_ch_session)

	// write buffer
	bufctrl := make(chan string)
	buf := NewBuffer(conn, bufctrl)
	go buf.Start()

	// cleanup work
	defer func() {
		close_work(&sess)
		close(timer_ch_session)
		close(sess.MQ)
		bufctrl <- "exit"
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

		case _ = <-timer_ch_session:
			if session_work(&sess, session_timeout) {
				return
			}
		}
	}
}
