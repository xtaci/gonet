package player

import "strings"
import "net"
import "encoding/binary"

// commands
func login(p string) string {
}

func echo(p string) string{
	return p
}

func attack(p string) string {
	return p
}

var header []byte

func send(conn net.Conn, p string) error {
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
	header = make([]byte,2)
	cmds := make(map[string]func(string)string)
	cmds["echo"] = echo

	for {
		msg := <-in
		if (msg == "CLIENTCLOSE") {
			break
		}

		params:= strings.SplitN(msg, " ", 2)

		cmd := cmds[params[0]]

		if cmd == nil {
			send(conn, "invalid command")
			continue;
		}

		switch len(params) {
		case 2:
			result := cmd(params[1])
			err := send(conn, result)
			if err != nil {
				conn.Close()
				break
			}
		}
	}
}
