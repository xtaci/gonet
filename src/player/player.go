package player

import "net"
import "time"
import "encoding/binary"
import . "types"
import . "db"

func send(conn net.Conn, p string) error {
	header := make([]byte,2)
	binary.BigEndian.PutUint16(header, uint16(len(p)));
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

func flush_timer(ud *User) {
	for {
		time.Sleep(10*time.Second)
		if ud.Id != 0 {
			DB.Flush(ud)
		}
		time.Sleep(4*time.Second)
	}
}

func NewPlayer(in chan string, conn net.Conn) {
	var user User
	user.MQ = make(chan string, 100)

	if send(conn, "Welcome") != nil {
		return
	}

	go flush_timer(&user)
L:
	for {
		select {
		case msg := <-in:
			if msg == "" {
				break L
			}

			result := exec_cli(&user, msg)

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

			result := exec_srv(&user, msg)

			if result != "" {
				err := send(conn, result)
				if err != nil {
					break L
				}
			}
		}
	}
}
