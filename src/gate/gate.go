package main

import "net"
import . "agent"
import . "db"
import "io"
import "os"
import "log"

func main() {
	config := read_config("./config.ini")
	if config["logfile"] != "" {
		f, err := os.OpenFile(config["logfile"], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

		if err != nil {
			log.Println("cannot open logfile %v\n", err)
			os.Exit(1)
		}
		var r Repeater
		r.Out1 = os.Stdout
		r.Out2 = f
		log.SetOutput(&r)
	}

	log.Println("Starting the server")
	StartDB(config)

	//	
	service := ":8080"
	if config["service"] != "" {
		service = config["service"]
	}

	//
	go SignalProc()

	log.Println("Service:", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}
		go handleClient(conn, config)
	}
}

func handleClient(conn net.Conn, config map[string]string) {
	defer conn.Close()

	header := make([]byte, 2)
	ch := make(chan []byte, 100)

	go StartAgent(ch, conn, config)

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
		size := int(header[0]<<8 | header[1])
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
		log.Println("Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
