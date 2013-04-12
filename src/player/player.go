package player

import "net"
import "time"
import "encoding/binary"
import . "types"
import . "db"
import "strconv"
import "cmd"

func send(conn net.Conn, p string) error {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, uint16(len(p)))
	_, err := conn.Write(header)
	if err != nil {
		println("Error send reply header:", err.Error())
		return err
	}

	_, err = conn.Write([]byte(p))
	if err != nil {
		println("Error send reply msg:", err.Error())
		return err
	}

	return nil
}

func timer_work(ud *User) {
	if ud.Id != 0 {
		DB.Flush(ud)
	}
}

func _timer(interval int, ch chan string) {
	defer func() {
		recover()
	}()

	func(ch chan string) {
		for {
			time.Sleep(time.Duration(interval) * time.Second)
			ch <- "timer"
		}
	}(ch)
}

func NewPlayer(in chan string, conn net.Conn, config map[string]string) {
	var user User
	user.MQ = make(chan string, 100)

	if send(conn, "Welcome") != nil {
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
		case msg := <-in:
			if msg == "" {
				break L
			}

			result := cmd.ExecCli(&user, msg)

			if result != "" {
				err := send(conn, result)
				if err != nil {
					break L
				}
			}

		case msg := <-user.MQ:
			if msg == "" {
				break L
			}

			result := cmd.ExecSrv(&user, msg)

			if result != "" {
				err := send(conn, result)
				if err != nil {
					break L
				}
			}
		case _ = <-timer_ch:
			timer_work(&user)
		}
	}

	close(timer_ch)
}
