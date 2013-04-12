package main

import "net"
import . "player"
import . "db"
import "io"
import "os"
import "fmt"

func main() {
	println("Starting the server")

	config := read_config("./config.ini")


	StartDB(config)

	//	
	service := ":8080"
	if config["service"] != "" {
		service = config["service"]
	}

	println("Service:", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	InitNames()
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
	ch := make(chan string, 100)

	go NewPlayer(ch, conn, config)

	for {
		// header
		n, err := io.ReadFull(conn, header)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			println("error receving header:", err)
			break
		}

		// data
		size := int(header[0]<<8 | header[1])
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err != nil {
			println("error receving msg:", err)
			break
		}
		ch <- string(data)
	}

	close(ch)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
