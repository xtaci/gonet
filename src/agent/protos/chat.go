package protos

import "strings"
import "strconv"
import . "types"
import "packet"
import "names"
import "log"

func Chat(ud *User, reader *packet.Packet) (ret []byte, err error) {
	user_id, err := reader.ReadString()
	if err !=nil {
		log.Println("Talk","read user_id failed.")
		return nil, err
	}

	msg, err := reader.ReadString()
	if err !=nil {
		log.Println("Talk","read msg failed.")
		return nil, err
	}

	id, _ := strconv.Atoi(user_id)
	ch := names.Query(id)
	if ch != nil {
		msg := []string{"MESG", string(ud.Id), msg}
		ch <- strings.Join(msg, " ")
	}

	writer := packet.PacketWriter()
	writer.WriteString("OK")
	return writer.Data(), err
}


