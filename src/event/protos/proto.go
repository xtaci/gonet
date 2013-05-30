package protos

import "misc/packet"

type ADD_REQ struct {
	F_tblname string
	F_oid     uint32
	F_user_id int32
	F_timeout int64
}

type CANCEL_REQ struct {
	F_event_id uint32
}

type INT struct {
	F_v uint32
}

func PKT_ADD_REQ(reader *packet.Packet) (tbl ADD_REQ, err error) {
	tbl.F_tblname, err = reader.ReadString()
	checkErr(err)

	tbl.F_oid, err = reader.ReadU32()
	checkErr(err)

	tbl.F_user_id, err = reader.ReadS32()
	checkErr(err)

	tbl.F_timeout, err = reader.ReadS64()
	checkErr(err)

	return
}

func PKT_CANCEL_REQ(reader *packet.Packet) (tbl CANCEL_REQ, err error) {
	tbl.F_event_id, err = reader.ReadU32()
	checkErr(err)

	return
}

func PKT_INT(reader *packet.Packet) (tbl INT, err error) {
	tbl.F_v, err = reader.ReadU32()
	checkErr(err)

	return
}
