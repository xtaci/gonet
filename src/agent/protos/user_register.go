package protos

import . "types"
import "misc/packet"

func UserRegister(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	udid, err := reader.ReadString()
	checkErr(err)
	name, err := reader.ReadString()
	checkErr(err)
	sex, err := reader.ReadU32()
	checkErr(err)

	println(udid, name, sex)
	return nil, nil
}
