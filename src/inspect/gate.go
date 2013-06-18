package inspect

import (
	"log"
	"net"
	"os"
)

import (
	"cfg"
)

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
		conn, err := listener.Accept()

		if err != nil {
			continue
		}
		handleClient(conn)
	}
}

var _conn net.Conn

func handleClient(conn net.Conn) {
	defer func() {
		conn.Close()
		recover()
	}()

	_conn = conn
	lex := NewLexer(_conn)
	yyParse(lex)
}

func checkError(err error) {
	if err != nil {
		log.Println("Fatal error: %v", err)
		os.Exit(-1)
	}
}
