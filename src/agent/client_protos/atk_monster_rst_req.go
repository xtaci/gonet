package protos

import . "types"
import "misc/packet"

func P_atk_monster_rst_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := PKT_atk_monster_rst_req(reader)
	println(tbl.F_protect_time)
	return
}
