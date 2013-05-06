package agent

import (
	"agent/ipc"
	"agent/protos"
	"misc/packet"
	. "types"
)

import (
	"fmt"
	"log"
	"runtime"
)

func UserRequestProxy(sess *Session, p []byte) []byte {
	defer _ProxyError()

	reader := packet.Reader(p)
	b, err := reader.ReadU16()

	if err != nil {
		log.Println("read protocol error")
	}

	handle := protos.ProtoHandler[b]
	if handle != nil {
		ret, err := handle(sess, reader)
		fmt.Println(ret)
		if err == nil {
			return ret
		}
	} else {
		log.Printf("no such protocol '%v'\n", b)
	}

	return nil
}

func IPCRequestProxy(sess *Session, p interface{}) []byte {
	defer _ProxyError()

	msg := p.(ipc.RequestType)
	handle := ipc.RequestHandler[msg.Code]
	if handle != nil {
		ipc, client := handle(sess, msg.Params)
		msg.CH <- ipc
		return client
	}

	return nil
}

func _ProxyError() {
	if x := recover(); x != nil {
		log.Printf("run time panic when processing user request: %v", x)
		for i := 0; i < 10; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				log.Printf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			}
		}
	}
}
