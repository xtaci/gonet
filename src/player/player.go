package player

import "net"
import "encoding/binary"

type UserData struct {
	id int;
	name string;
	mq chan string;
}

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

func Start(in chan string, conn net.Conn) {
	var user UserData
	user.mq = make(chan string, 100)

	if send(conn, "Welcome") != nil {
		return
	}
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
					break
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
					break
				}
			}
		}
	}
}
