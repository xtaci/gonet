package protos

import "misc/packet"

type ADD_REQ struct {
	F_user_id int32
	F_obj_id int32
	F_obj_type int32
	F_obj_nextlevel int32
	F_timeout int64
}

type CANCEL_REQ struct {
	F_event_id int32
}

type LONG struct {
	F_v int64
}

type INT struct {
	F_v int32
}

func PKT_ADD_REQ(reader *packet.Packet)(tbl ADD_REQ, err error){
	tbl.F_user_id,err = reader.ReadS32()
	checkErr(err)
	tbl.F_obj_id,err = reader.ReadS32()
	checkErr(err)
	tbl.F_obj_type,err = reader.ReadS32()
	checkErr(err)
	tbl.F_obj_nextlevel,err = reader.ReadS32()
	checkErr(err)
	return
}

func PKT_CANCEL_REQ(reader *packet.Packet)(tbl CANCEL_REQ, err error){
	tbl.F_event_id,err = reader.ReadS32()
	checkErr(err)
	return
}

func PKT_LONG(reader *packet.Packet)(tbl LONG, err error){
	return
}

func PKT_INT(reader *packet.Packet)(tbl INT, err error){
	tbl.F_v,err = reader.ReadS32()
	checkErr(err)
	return
}

