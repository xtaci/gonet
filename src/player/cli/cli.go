package cli

import "strings"
import "strconv"
import "db/user"
import "db/city"
import . "types"
import "names"
import "log"
import "packet"

// commands from client
func Login(ud *User, reader *packet.Packet) (ret []byte, err error) {
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

func Echo(ud *User, reader *packet.Packet) (ret []byte, err error) {
	msg, err := reader.ReadString()
	if err !=nil {
		log.Println("Echo","read message failed.")
		return nil, err
	}

	writer := packet.PacketWriter()
	writer.WriteString(msg)
	return writer.Data(), err
}

func Talk(ud *User, reader *packet.Packet) (ret []byte, err error) {
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

/*
func Newcity(ud *User, reader *Packet) (ret []byte, err error) {
	newcity := City { Name:p, OwnerId:ud.Id }
	ud.Cities = append(ud.Cities, newcity)
	city.Create(&ud.Cities[len(ud.Cities)-1])
	writer := PacketWriter()
	writer.WriteString("OK")
	return writer.Data(), err
}
*/
