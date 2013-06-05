package protos

import (
	"fmt"
	"log"
	"runtime"
)

import (
	. "helper"
	"misc/packet"
)

//--------------------------------------------------------- send
func _send(seqid uint64, data []byte, output chan []byte) {
	writer := packet.Writer()
	writer.WriteU16(uint16(len(data)) + 8)
	writer.WriteU64(seqid) // piggyback seq id
	writer.WriteRawBytes(data)
	output <- writer.Data()
}

func HandleRequest(reader *packet.Packet, output chan []byte) {
	defer PrintPanicStack()

	seqid, err := reader.ReadU64() // read seqid
	if err != nil {
		log.Println("Read Sequence Id failed.", err)
		return
	}

	b, err := reader.ReadU16()
	if err != nil {
		log.Println("read protocol error")
		return
	}

	fmt.Println("proto: ", b)
	handle := ProtoHandler[b]
	if handle != nil {
		ret := handle(reader)
		if len(ret) != 0 {
			_send(seqid, ret, output)
		}
	}
}

func P_ping_req(reader *packet.Packet) []byte {
	tbl, _ := PKT_INT(reader)
	ret := INT{tbl.F_v}
	return packet.Pack(Code["ping_ack"], ret, nil)
}

func P_add_req(reader *packet.Packet) []byte {
	tbl, _ := PKT_ADD_REQ(reader)
	ret := INT{0}

	fmt.Println(tbl)
	return packet.Pack(Code["add_ack"], ret, nil)
}

func checkErr(err error) {
	if err != nil {
		funcName, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("ERR:%v,[func:%v,file:%v,line:%v]\n", err, runtime.FuncForPC(funcName).Name(), file, line)
		}

		panic("error occured in Stats Protocol Module")
	}
}
