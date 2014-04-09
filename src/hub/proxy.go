package main

import (
	"log"
	"time"
)

import (
	"helper"
	"hub/protos"
	"misc/packet"
)

//------------------------------------------------ Game Server Request Proxy
func GSProxy(hostid int32, reader *packet.Packet) (ret []byte) {
	defer helper.PrintPanicStack()

	// read protocol number
	b, err := reader.ReadS16()
	if err != nil {
		log.Println("read protocol error")
		return
	}

	// get handler
	handle := protos.ProtoHandler[b]
	if handle == nil {
		log.Println("service not bind", b)
		return
	}

	// call handler
	start := time.Now()
	ret = handle(hostid, reader)
	end := time.Now()
	log.Printf("code: %v %v TIME:%v\n", b, protos.RCode[b], end.Sub(start))

	return ret
}
