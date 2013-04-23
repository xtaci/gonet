package protos

import . "types"
import "misc/packet"
import "hub/names"
import "db/user"

func _user_login_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {

	tbl,_ := pktread_user_login_info(reader)

	if tbl.new_user == 0 {
		if user.LoginMAC(sess.User.Mac, &sess.User) {
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
