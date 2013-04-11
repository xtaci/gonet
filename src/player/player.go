package player

import "net"
import "time"
import "encoding/binary"

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

func (user *User) flush_timer() {
	for {
		time.Sleep(10*time.Second)
		if user.id != 0 {
			DB.Flush(user)
		}
		time.Sleep(4*time.Second)
	}
}

func NewPlayer(in chan string, conn net.Conn) {
	var user User
	user.mq = make(chan string, 100)

	if send(conn, "Welcome") != nil {
		return
	}

	go user.flush_timer()
L:
	for {
		select {
		case msg := <-in:
			if msg == "" {
				break L
			}

			result := user.exec_cli(msg)

			if result != "" {
				err := send(conn, result)
				if err != nil {
					break L
				}
			}

		case msg := <-user.mq:
			if msg == "" {
				break L
			}

			result := user.exec_srv(msg)

			if result != "" {
				err := send(conn, result)
				if err != nil {
					break L
				}
			}
		}
	}
}
