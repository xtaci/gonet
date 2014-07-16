package main

import (
	"log"
)

import (
	. "helper"
	"misc/packet"
	"stats/protos"
)

func HandleRequest(reader *packet.Packet) {
	defer PrintPanicStack()
	b, err := reader.ReadS16()
	if err != nil {
		log.Println("read protocol error")
		return
	}

	handle := protos.ProtoHandler[b]
	//DEBUG("=== stats protocal====", b)
	if handle == nil {
		log.Println("service not bind", b)
		return
	}
	handle(reader)
}
