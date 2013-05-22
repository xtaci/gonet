package main

import (
	"log"
	"net"
)

import (
	"misc/packet"
	"cooldown/protos"
)

const (
	MAXCHAN = 200000
)

func init() {
	log.SetPrefix("[CD]")
}

//--------------------------------------------------------- send
func _send(seqid uint64, data []byte, output chan []byte) {
	writer := packet.Writer()
	writer.WriteU16(uint16(len(data)) + 8)
	writer.WriteU64(seqid) // piggyback seq id
	writer.WriteRawBytes(data)
	output <- writer.Data()
}

//------------------------------------------------ CoolDown Agent
func CDAgent(incoming chan []byte, conn net.Conn) {
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
			go protos.HandleRequest(reader, output)
		}
	}
}

//----------------------------------------------- write buffer 
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
