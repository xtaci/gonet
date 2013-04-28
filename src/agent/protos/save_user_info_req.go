package protos

import (
	. "types"
	"misc/packet"
	"db/user_tbl"
)

func _save_user_info_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl,_ := pktread_user_archives_info(reader)
	sess.User.Archives = tbl.F_archives
	user_tbl.Flush(&sess.User)
	return
}
