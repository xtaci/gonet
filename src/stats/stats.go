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
	"cfg"
	"helper"
)

const (
	DEFAULT_SERVICE = ":8891"
)

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

//----------------------------------------------- Stats Server start
func StatsStart() {
	config := cfg.Get()
	if config["profile"] == "true" {
		helper.SetMemProfileRate(1)
		defer func() {
			helper.GC()
			helper.DumpHeap()
			helper.PrintGCSummary()
		}()
	}

	// start logger
	if config["stats_log"] != "" {
		cfg.StartLogger(config["stats_log"])
	}

	log.Println("Starting Stats Server")
	go SignalProc()

	// Listen
	service := DEFAULT_SERVICE
	if config["stats_service"] != "" {
		service = config["stats_service"]
	}

	log.Println("Stats Service:", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	log.Println("Stats Server OK.")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		helper.SetConnParam(conn)
		go handleClient(conn)
	}
}

//----------------------------------------------- handle cooldown request
func handleClient(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, 2)
	ch := make(chan []byte, 8192)

	go StatsAgent(ch, conn)

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

func main() {
	StatsStart()
}
