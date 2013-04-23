package agent

import "net"
import "time"
import "encoding/binary"
import . "types"
import "db/user"
import "db/city"
import "strconv"
import "names"
import "log"
import "fmt"

func send(conn net.Conn, p []byte) error {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, uint16(len(p)))
	_, err := conn.Write(header)
	if err != nil {
		log.Println("Error send reply header:", err.Error())
		return err
	}

	_, err = conn.Write(p)
	if err != nil {
		log.Println("Error send reply msg:", err.Error())
		return err
	}

	return nil
}

func timer_work(sess *Session) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic when flushing database: %v", x)
		}
	}()

	if sess.User.Id!= 0 {
		user.Flush(&sess.User)
		for i := range sess.Cities {
			city.Flush(&sess.Cities[i])
		}
	}
}

func _timer(interval int, ch chan string) {
	defer func() {
		recover()
	}()

	for {
		time.Sleep(time.Duration(interval) * time.Second)
		ch <- "ding!dong!"
	}
}

func StartAgent(in chan []byte, conn net.Conn, config map[string]string) {
	var sess Session
	sess.MQ = make(chan string, 128)
	timer_ch := make(chan string)

	flush_interval := 300 // sec
	if config["flush_interval"] != "" {
		flush_interval,_ = strconv.Atoi(config["flush_interval"])
	}

	go _timer(flush_interval, timer_ch)
L:
	for {
		select {
		case msg,ok := <-in:
			if !ok {
				break L
			}

			if result := ExecCli(&sess, msg); result != nil {
				fmt.Println(result)
				err := send(conn, result)
				if err != nil {
					break L
				}
			}

		case msg,ok := <-sess.MQ:
			if !ok {
				break L
			}

			result := ExecSrv(&sess, msg)

			if result != "" {
				err := send(conn, []byte(result))
				if err != nil {
					break L
				}
			}
		case _ = <-timer_ch:
			timer_work(&sess)
		}
	}

	// cleanup
	names.Unregister(sess.User.Id)
	close(timer_ch)
}
