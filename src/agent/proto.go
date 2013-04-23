package agent

import "strings"
import . "types"
import srv "agent/srv"
import "agent/protos"
import "log"
import "misc/packet"

func ExecCli(sess *Session, p []byte) []byte {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic when processing user request: %v", x)
		}
	}()

	reader := packet.PacketReader(p)

	b, err := reader.ReadU16()

	if err != nil {
		log.Println("read protocol error")
	}

	handle := ProtoHandler[b]
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

func ExecSrv(sess *Session, msg string) string {
	params := strings.SplitN(msg, " ", 2)
	switch params[0] {
	case "MESG":
		return srv.Mesg(&sess.User, params[1])
	case "ATTACKED":
		return srv.Attacked(&sess.User, params[1])
	}

	return ""
}

//---------------------------------------------------------Handler Binding
var ProtoHandler map[uint16]func(*Session, *packet.Packet) ([]byte, error)

func init() {
	ProtoHandler = make(map[uint16]func(*Session, *packet.Packet) ([]byte, error))
	ProtoHandler[0] = protos.HeartBeat
	ProtoHandler['R'] = protos.UserRegister
	ProtoHandler[3] = protos.UserLogin
	ProtoHandler[9] = protos.Chat
	ProtoHandler[11] = protos.UserLogout
}
