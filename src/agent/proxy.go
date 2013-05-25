package agent

import (
	"agent/client_protos"
	"agent/ipc"
	"misc/packet"
	. "types"
)

import (
	"fmt"
	"log"
	"runtime"
)

//----------------------------------------------- client protocol handle proxy
func UserRequestProxy(sess *Session, p []byte) []byte {
	defer _ProxyError()

	reader := packet.Reader(p)
	b, err := reader.ReadU16()

	if err != nil {
		log.Println("read protocol error")
	}

	log.Printf("code:%v,user:%v\n", b, sess.Basic.Id)

	handle := protos.ProtoHandler[b]
	if handle != nil {
		ret, err := handle(sess, reader)
		fmt.Println(ret)
		if err == nil {
			return ret
		}
	}

	return nil
}

//----------------------------------------------- IPC proxy
func IPCRequestProxy(sess *Session, p *IPCObject) {
	defer _ProxyError()
	handle := ipc.IPCHandler[p.Service]
	log.Printf("ipc:%v,user:%v\n", p.Service, sess.Basic.Id)

	if handle != nil {
		handle(sess, p)
	}
}

func _ProxyError() {
	if x := recover(); x != nil {
		log.Printf("run time panic when processing request: %v", x)
		for i := 0; i < 10; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				log.Printf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			}
		}
	}
}
