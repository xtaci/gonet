package agent

import . "types"
import "agent/ipc"
import "agent/protos"
import "log"
import "misc/packet"

func UserRequestProxy(sess *Session, p []byte) []byte {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic when processing user request: %v", x)
		}
	}()

	reader := packet.Reader(p)

	b, err := reader.ReadU16()

	if err != nil {
		log.Println("read protocol error")
	}

	handle := protos.ProtoHandler[b]
	if handle != nil {
		ret, err := handle(sess, reader)

		if err == nil {
			return ret
		}
	} else {
		log.Printf("no such protocol '%v'\n", b)
	}

	return nil
}

func IPCRequestProxy(sess *Session, p interface{}) []byte {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic when processing IPC request: %v", x)
		}
	}()

	msg := p.(ipc.RequestType)
	handle := ipc.RequestHandler[msg.Code]
	if handle !=nil {
		msg.CH <- handle(sess, msg.Params)
	}

	return nil
}
