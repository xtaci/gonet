package protos

import . "types"
import "misc/packet"
import "names"
import "db/user"
import "db/city"

func UserLogin(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	name, err := reader.ReadString()
	checkErr(err)
	pass, err := reader.ReadString()
	checkErr(err)

	if user.Login(name, pass, &sess.User) {
		names.Register(sess.MQ, sess.User.Id)
		// load cities
		sess.Cities = city.LoadCities(sess.User.Id)
	}

	writer := packet.PacketWriter()
	writer.WriteString("true")

	return writer.Data(), err
}
