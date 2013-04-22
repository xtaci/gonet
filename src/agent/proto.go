package agent

import "strings"
import . "types"
import srv "agent/srv"
import "agent/protos"
import "log"
import "packet"

func ExecCli(ud *User, p []byte) []byte {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("run time panic when processing user request: %v", x)
		}
	}()

	reader := packet.PacketReader(p)

	b, err := reader.ReadByte()

	if err!=nil {
		log.Println("read protocol error")
	}

	handle := ProtoHandler[uint16(b)]
	if handle != nil {
		ret, err := handle(ud, reader)

		if err == nil {
			return ret
		}
	} else {
		log.Printf("no such protocol '%v'\n", b)
	}

	return nil
}

func ExecSrv(ud *User, msg string) string {
	params := strings.SplitN(msg, " ", 2)
	switch params[0] {
	case "MESG":
		return srv.Mesg(ud, params[1])
	case "ATTACKED":
		return srv.Attacked(ud, params[1])
	}

	return ""
}

//---------------------------------------------------------Handler Binding
var ProtoHandler map[uint16]func(*User, *packet.Packet)([]byte, error)
func init() {
	ProtoHandler = make(map[uint16]func(*User, *packet.Packet)([]byte, error))
	ProtoHandler['R'] = protos.UserRegister
	ProtoHandler[3] = protos.UserLogin
	ProtoHandler[9] = protos.Chat
	ProtoHandler[11] = protos.UserLogout
}
