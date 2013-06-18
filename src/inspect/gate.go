package inspect

import (
	"fmt"
	"log"
	"net"
	"os"
)

import (
	"cfg"
)

var conn net.Conn

func StartInspect() {
	config := cfg.Get()
	// Listen
	service := ":8800"
	if config["inspect"] != "" {
		service = config["inspect"]
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err = listener.Accept()

		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Fprintln(conn, x)
		}
		conn.Close()
	}()

	fmt.Fprintln(conn, "GameServer Console")
	fmt.Fprintln(conn, `type 'help' for usage`)
	prompt(conn)
	lex := NewLexer(conn)
	yyParse(lex)
}

func checkError(err error) {
	if err != nil {
		log.Println("Fatal error: %v", err)
		os.Exit(-1)
	}
}
