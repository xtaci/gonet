package main

import (
	"log"
	"net"
	"os"
	"strconv"
)

import (
	"cfg"
)

const (
	DEFAULT_SERVICE = ":8891"
)

//----------------------------------------------- Stats Server start
func main() {
	defer func() {
		if x := recover(); x != nil {
			log.Println("caught panic in main()", x)
		}
	}()

	config := cfg.Get()
	// start logger
	if config["stats_log"] != "" {
		cfg.StartLogger(config["stats_log"])
	}

	log.Println("Starting Stats Server")
	go SignalProc()
	go SysRoutine()

	// Listen
	service := DEFAULT_SERVICE
	if config["stats_service"] != "" {
		service = config["stats_service"]
	}

	log.Println("Stats Service:", service)
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	udpconn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	log.Println("Stats Server OK.")
	handleClient(udpconn)
}

//----------------------------------------------- handle cooldown request
func handleClient(conn *net.UDPConn) {
	// init receive buffer
	config := cfg.Get()
	maxchan, e := strconv.Atoi(config["stats_max_queue_size"])
	if e != nil {
		maxchan = DEFAULT_MAX_QUEUE
		log.Println("cannot parse stats_max_queue_size from config", e)
	}

	ch := make(chan []byte, maxchan)
	defer close(ch)

	go StatsAgent(ch, conn)

	// loop receiving
	for {
		// udp receive buffer, max 512 packet
		data := make([]byte, 512)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println("read udp failed", n, addr, err)
			continue
		}

		ch <- data[:n]
	}
}

func checkError(err error) {
	if err != nil {
		log.Printf("Fatal error: %v", err)
		os.Exit(-1)
	}
}
