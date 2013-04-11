package player

import "strings"
import "net"
import "encoding/binary"

var mq chan string

type UserData struct {
	id int;
	name string;
	mq chan string;
}

var user UserData

// commands from server
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
	user.mq = make(chan string)

	client_cmds := map[string]func(string) string  {
		"echo":Client_echo,
		"login":Client_login,
		"talk":Client_talk,
		"attack":Client_attack,
	}

	server_cmds := map[string]func(string) string  {
		"mesg":Server_mesg,
		"attackedby": Server_attackedby,
	}

	for {
		select {
		case msg := <-in:
			if msg == "CLIENTCLOSE" || msg == "" {
				break
			}

			params:= strings.SplitN(msg, " ", 2)

			cmd := client_cmds[params[0]]

			if cmd == nil {
				send(conn, "invalid command")
				continue;
			}

			result := cmd(params[1])
			err := send(conn, result)
			if err != nil {
				conn.Close()
				break
			}

		case msg := <-mq:
			if msg == "" {
				break
			}

			params:= strings.SplitN(msg, " ", 2)

			cmd := server_cmds[params[0]]

			if cmd == nil { continue; }

			result := cmd(params[1])
			err := send(conn, result)
			if err != nil {
				conn.Close()
				break
			}
		}
	}
}
