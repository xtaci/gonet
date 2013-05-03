package protos

import (
	"db/user_tbl"
	"misc/packet"
	. "types"
)

func _save_user_info_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := pktread_user_archives_info(reader)
	sess.User.Archives = tbl.F_archives
	user_tbl.Sync(&sess.User)
	return
}
