package protos

import . "types"
import "misc/packet"

func P_atk_player_rst_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	tbl, _ := PKT_atk_player_rst_req(reader)
	println(tbl.F_rst)
	return
}
