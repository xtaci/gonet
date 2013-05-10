package agent

import (
	"agent/ipc"
	"agent/protos"
	"misc/packet"
	. "types"
	"cfg"
)

import (
	"fmt"
	"log"
	"runtime"
	"os"
)

var (
	proto_logger *log.Logger
	ipc_logger *log.Logger
)

//----------------------------------------------- client protocol handle proxy
func UserRequestProxy(sess *Session, p []byte) []byte {
	defer _ProxyError()

	reader := packet.Reader(p)
	b, err := reader.ReadU16()

	if err != nil {
		log.Println("read protocol error")
	}

	proto_logger.Printf("tos:%v,user:%v\n", b, sess.User.Id)

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
func IPCRequestProxy(sess *Session, p interface{}) []byte {
	defer _ProxyError()
	msg := p.(ipc.RequestType)
	handle := ipc.RequestHandler[msg.Code]
	ipc_logger.Printf("ipc:%v,user:%v\n", msg.Code, sess.User.Id)

	if handle != nil {
		ipc, client := handle(sess, msg.Params)
		msg.CH <- ipc
		return client
	}

	return nil
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

func init() {
	config := cfg.Get()
	proto_logfile, err := os.OpenFile(config["proto_logfile"], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Println("cannot open proto logfile %v\n", err)
		os.Exit(1)
	}

	proto_logger = log.New(proto_logfile, "", log.LstdFlags)

	//
	ipc_logfile, err := os.OpenFile(config["ipc_logfile"], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Println("cannot open ipc logfile %v\n", err)
		os.Exit(1)
	}

	ipc_logger = log.New(ipc_logfile, "", log.LstdFlags)
}
