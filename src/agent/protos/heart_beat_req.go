package protos

import . "types"
import "misc/packet"
import "time"

func P_heart_beat_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	sess.HeartBeat = time.Now()
	return
}
