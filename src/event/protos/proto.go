package protos

import "misc/packet"

type ADD_EVENT struct {
	F_type    int16
	F_user_id int32
	F_timeout int64
	F_params  []byte
}

type CANCEL_EVENT struct {
	F_event_id int32
}

type EVENT_ID struct {
	F_event_id int32
}

type INT struct {
	F_v int32
}

func PKT_ADD_EVENT(reader *packet.Packet) (tbl ADD_EVENT, err error) {
	tbl.F_type, err = reader.ReadS16()
	checkErr(err)

	tbl.F_user_id, err = reader.ReadS32()
	checkErr(err)

	tbl.F_timeout, err = reader.ReadS64()
	checkErr(err)

	tbl.F_params, err = reader.ReadBytes()
	checkErr(err)

	return
}

func PKT_CANCEL_EVENT(reader *packet.Packet) (tbl CANCEL_EVENT, err error) {
	tbl.F_event_id, err = reader.ReadS32()
	checkErr(err)

	return
}

func PKT_EVENT_ID(reader *packet.Packet) (tbl EVENT_ID, err error) {
	tbl.F_event_id, err = reader.ReadS32()
	checkErr(err)

	return
}

func PKT_INT(reader *packet.Packet) (tbl INT, err error) {
	tbl.F_v, err = reader.ReadS32()
	checkErr(err)

	return
}
