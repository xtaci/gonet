package main

import (
	"encoding/binary"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"time"
)

import (
	"cfg"
	. "helper"
	"misc/geoip"
	. "types"
)

func main() {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in main()", x)
		}
	}()

	go func() {
		INFO(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	// start basic services
	startup()

	// Listen
	config := cfg.Get()
	service := ":8080"
	if config["service"] != "" {
		service = config["service"]
	}

	INFO("Service:", service)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	INFO("Game Server OK.")

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			WARN("accept failed", err)
			continue
		}

		go handleClient(conn)
	}
}

//----------------------------------------------- start a goroutine when a new connection is accepted
func handleClient(conn *net.TCPConn) {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in handleClient", x)
		}
	}()

	// input buffer
	config := cfg.Get()
	inqueue_size, err := strconv.Atoi(config["inqueue_size"])
	if err != nil {
		inqueue_size = DEFAULT_INQUEUE_SIZE
		WARN("cannot parse inqueue_size from config", err, "using default:", inqueue_size)
	}

	// init
	header := make([]byte, 2)
	in := make(chan []byte, inqueue_size)
	bufctrl := make(chan bool)

	defer func() {
		close(bufctrl)
		close(in)
	}()

	// create new session
	var sess Session
	sess.IP = net.ParseIP(strings.Split(conn.RemoteAddr().String(), ":")[0])
	NOTICE("new connection from:", sess.IP, "country:", geoip.Query(sess.IP))

	// create write buffer
	out := NewBuffer(&sess, conn, bufctrl)
	go out.Start()

	// start agent!!
	go StartAgent(&sess, in, out)

	for {
		// header
		conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		n, err := io.ReadFull(conn, header)
		if err != nil {
			WARN("error receiving header, bytes:", n, "reason:", err)
			break
		}

		// data
		size := binary.BigEndian.Uint16(header)
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)
		if err != nil {
			WARN("error receiving msg, bytes:", n, "reason:", err)
			break
		}

		// NOTICE: slice is passed by reference; don't re-use a single buffer.
		select {
		case in <- data:
		case <-time.After(MAX_DELAY_IN * time.Second):
			WARN("server busy or agent closed, session flag:", sess.Flag)
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		ERR("Fatal error:", err)
		os.Exit(-1)
	}
}
