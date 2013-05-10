package main

import (
	"net"
	"log"
)

import (
	"misc/packet"
	"hub/protos"
)

//--------------------------------------------------------- Hub processing
func HubAgent(in chan []byte, conn net.Conn) {
	for {
		msg, ok := <-in

		if !ok {
			return
		}

		if result := protos.HandleRequest(msg); result != nil {
			headwriter := packet.Writer()
			headwriter.WriteU16(uint16(len(result)))

			_, err := conn.Write(headwriter.Data())
			if err != nil {
				log.Println("Error send reply header:", err.Error())
				return
			}

			_, err = conn.Write(result)
			if err != nil {
				log.Println("Error send reply msg:", err.Error())
				return
			}
		}
	}
}

func init() {
	log.SetPrefix("[HUB]")
}
