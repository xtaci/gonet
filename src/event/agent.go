package main

import (
	"log"
	"net"
)

import (
	"event/protos"
	"misc/packet"
)

const (
	MAXCHAN = 1000000
)

func init() {
	log.SetPrefix("[EVENT]")
}

//----------------------------------------------------------- Event Server Agent
func EventAgent(incoming chan []byte, conn net.Conn) {
	// output buffer
	output := make(chan []byte, MAXCHAN)
	go _write_routine(output, conn)

	defer func() {
		close(output)
	}()

	for {
		select {
		case msg, ok := <-incoming:
			if !ok {
				return
			}

			reader := packet.Reader(msg)
			protos.HandleRequest(reader, output)
		}
	}
}

//---------------------------------------------------------- write buffer
func _write_routine(output chan []byte, conn net.Conn) {
	for {
		msg, ok := <-output
		if !ok {
			break
		}

		_, err := conn.Write(msg) // write operation is assumed to be atomic
		if err != nil {
			log.Println("Error send reply to GS:", err)
		}
	}
}
