package main

import (
	"encoding/binary"
	"flag"
	"io"
	"log"
	"net"
	"os"
)

import (
	"agent/event_client"
	"agent/hub_client"
	"agent/stats_client"
	"cfg"
	"inspect"
)

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

//----------------------------------------------- Game Server Start
func main() {
	// start logger
	config := cfg.Get()
	if config["gs_log"] != "" {
		cfg.StartLogger(config["gs_log"])
	}

	// inspector
	go inspect.StartInspect()

	// dial HUB
	hub_client.DialHub()
	event_client.DialEvent()
	stats_client.DialStats()

	log.Println("Starting the server.")
	// signal
	go SignalProc()

	// Listen
	service := ":8080"
	if config["service"] != "" {
		service = config["service"]
	}

	log.Println("Service:", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	log.Println("Game Server OK.")

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		// DoS prevention
		IP := net.ParseIP(conn.RemoteAddr().String())
		if !IsBanned(IP) {
			go handleClient(conn)
		} else {
			conn.Close()
		}
	}
}

//----------------------------------------------- start a goroutine when a new connection is accepted
func handleClient(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, 2)
	ch := make(chan []byte, 10)

	go StartAgent(ch, conn)

	for {
		// header
		n, err := io.ReadFull(conn, header)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Println("error receiving header:", err)
			break
		}

		// data
		size := binary.BigEndian.Uint16(header)
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err != nil {
			log.Println("error receiving msg:", err)
			break
		}
		ch <- data
	}

	close(ch)
}

func checkError(err error) {
	if err != nil {
		log.Println("Fatal error: %v", err)
		os.Exit(-1)
	}
}
