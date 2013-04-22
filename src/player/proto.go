package player

import "strings"
import . "types"
import srv "player/srv"
import "player/protos"
import "log"
import "packet"

func ExecCli(ud *User, p []byte) []byte {
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

	return []byte{0}
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
	ProtoHandler[1] = protos.UserRegister
	ProtoHandler[3] = protos.UserLogin
	ProtoHandler[9] = protos.Chat
	ProtoHandler[11] = protos.UserLogout
}
