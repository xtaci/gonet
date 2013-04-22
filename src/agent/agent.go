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

func timer_work(ud *User) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic when flushing database: %v", x)
		}
	}()

	if ud.Id != 0 {
		_flush_all(ud)
	}
}

func _flush_all(ud *User) {
	user.Flush(ud)
	for i := range ud.Cities {
		city.Flush(&ud.Cities[i])
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
	var user User
	user.MQ = make(chan string, 100)

	if send(conn, []byte("Welcome")) != nil {
		return
	}

	timer_ch := make(chan string, 10)

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

			if result := ExecCli(&user, msg); result != nil {
				err := send(conn, result)
				if err != nil {
					break L
				}
			}

		case msg,ok := <-user.MQ:
			if !ok {
				break L
			}

			result := ExecSrv(&user, msg)

			if result != "" {
				err := send(conn, []byte(result))
				if err != nil {
					break L
				}
			}
		case _ = <-timer_ch:
			timer_work(&user)
		}
	}

	// cleanup
	names.Unregister(user.Id)
	close(timer_ch)
}
