package main

import (
	"io"
	"log"
	"net"
	"os"
)

import (
	"cfg"
)

const (
	DEFAULT_SERVICE = ":8890"
)

//----------------------------------------------- Event Server start
func EventStart() {
	// start logger
	config := cfg.Get()
	if config["event_log"] != "" {
		cfg.StartLogger(config["event_log"])
	}

	log.Println("Starting Event Server")
	startup_work()
	go SignalProc()

	// Listen
	service := DEFAULT_SERVICE
	if config["event_service"] != "" {
		service = config["event_service"]
	}

	log.Println("Event Service:", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	log.Println("Event Server OK.")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		conn.SetNoDelay(false)
		go handleClient(conn)
	}
}

//----------------------------------------------- handle cooldown request
func handleClient(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, 2)
	ch := make(chan []byte, 8192)

	go EventAgent(ch, conn)

	for {
		// header
		n, err := io.ReadFull(conn, header)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Println("error receving header:", err)
			break
		}

		// data
		size := int(header[0])<<8 | int(header[1])
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err != nil {
			log.Println("error receving msg:", err)
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

func main() {
	EventStart()
}
