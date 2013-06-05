package protos

import "misc/packet"

type ADD_REQ struct {
	F_data []byte
}

type INT struct {
	F_v int32
}

func PKT_ADD_REQ(reader *packet.Packet) (tbl ADD_REQ, err error) {
	tbl.F_data, err = reader.ReadBytes()
	checkErr(err)

	return
}

func PKT_INT(reader *packet.Packet) (tbl INT, err error) {
	tbl.F_v, err = reader.ReadS32()
	checkErr(err)

	return
}
