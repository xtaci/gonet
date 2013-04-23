package protos

import . "types"
import "misc/packet"
import "time"

func HeartBeat(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	sess.HeartBeat = time.Now()
	return
}
