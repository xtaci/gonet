package protos

import "strings"
import "strconv"
import . "types"
import "misc/packet"
import "hub/names"

func Chat(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	user_id, err := reader.ReadString()
	checkErr(err)
	msg, err := reader.ReadString()
	checkErr(err)
	id, _ := strconv.Atoi(user_id)
	ch := names.Query(id)
	if ch != nil {
		msg := []string{"MESG", string(sess.User.Id), msg}
		ch <- strings.Join(msg, " ")
	}

	writer := packet.PacketWriter()
	writer.WriteString("OK")
	return writer.Data(), err
}
