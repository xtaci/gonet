package protos

import . "types"
import "packet"
import "log"
import "names"
import "db/user"
import "db/city"

func UserLogin(ud *User, reader *packet.Packet) (ret []byte, err error) {
	name, err := reader.ReadString()
	if err !=nil {
		log.Println("Login", "read name failed.")
		return nil, err
	}

	pass, err := reader.ReadString()
	if err !=nil {
		log.Println("Login","read pass failed.")
		return nil, err
	}

	if user.Login(name, pass, ud) {
		names.Register(ud.MQ, ud.Id)
		// load cities
		ud.Cities = city.LoadCities(ud.Id)
	}

	writer := packet.PacketWriter()
	writer.WriteString("true")

	return writer.Data(), err
}


