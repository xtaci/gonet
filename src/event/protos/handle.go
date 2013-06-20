package protos

import (
	"fmt"
	"log"
	"runtime"
)

import (
	"event/core"
	. "helper"
	"misc/packet"
)

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
			SendChan(seqid, ret, output)
		}
	}
}

func P_ping_req(reader *packet.Packet) []byte {
	tbl, _ := PKT_INT(reader)
	ret := INT{tbl.F_v}
	return packet.Pack(-1, &ret, nil)
}

func P_add_req(reader *packet.Packet) []byte {
	tbl, _ := PKT_ADD_EVENT(reader)
	event_id := core.Add(tbl.F_type, tbl.F_user_id, tbl.F_timeout, tbl.F_params)
	ret := INT{event_id}

	return packet.Pack(-1, &ret, nil)
}

func P_cancel_req(reader *packet.Packet) []byte {
	tbl, _ := PKT_CANCEL_EVENT(reader)
	core.Cancel(tbl.F_event_id)
	ret := INT{1}

	return packet.Pack(-1, &ret, nil)
}

func checkErr(err error) {
	if err != nil {
		funcName, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("ERR:%v,[func:%v,file:%v,line:%v]\n", err, runtime.FuncForPC(funcName).Name(), file, line)
		}

		panic("error occured in Event Protocol Module")
	}
}
