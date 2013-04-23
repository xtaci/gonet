package protos

import . "types"
import "misc/packet"
import "hub/names"
import "db/user"
import "time"

func UserLogin(sess *Session, reader *packet.Packet) (ret []byte, err error) {

	tbl,_ := pktread_user_login_info(reader)

	if tbl.new_user == 0 {
		if user.LoginMAC(mac, &sess.User) {
			names.Register(sess.MQ, sess.User.Id)
		}
	} else {
		sess.User.Name = tbl.user_name
		sess.User.Mac = tbl.mac_addr

		if user.New(&sess.User) {
			names.Register(sess.MQ, sess.User.Id)
		}
	}

	return nil,nil
}
