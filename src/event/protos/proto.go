package protos

import "misc/packet"

type ADD_REQ struct {
	F_oid int32
	F_user_id int32
	F_timeout int64
}

type ADD_MOVES_REQ struct {
	F_user_id int32
	F_moves []MOVE
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

type MOVE struct {
	F_OID int32
	F_X int32
	F_Y int32
}

func PKT_ADD_REQ(reader *packet.Packet)(tbl ADD_REQ, err error){
	tbl.F_oid,err = reader.ReadS32()
	checkErr(err)

	tbl.F_user_id,err = reader.ReadS32()
	checkErr(err)

	tbl.F_timeout,err = reader.ReadS64()
	checkErr(err)

	return
}

func PKT_ADD_MOVES_REQ(reader *packet.Packet)(tbl ADD_MOVES_REQ, err error){
	tbl.F_user_id,err = reader.ReadS32()
	checkErr(err)

	narr := uint16(0)

	narr,err = reader.ReadU16()
	checkErr(err)

	tbl.F_moves=make([]MOVE,narr)
	for i:=0;i<int(narr);i++ {
		tbl.F_moves[i], err = PKT_MOVE(reader)
	}

	return
}

func PKT_CANCEL_REQ(reader *packet.Packet)(tbl CANCEL_REQ, err error){
	tbl.F_event_id,err = reader.ReadS32()
	checkErr(err)

	return
}

func PKT_LONG(reader *packet.Packet)(tbl LONG, err error){
	tbl.F_v,err = reader.ReadS64()
	checkErr(err)

	return
}

func PKT_INT(reader *packet.Packet)(tbl INT, err error){
	tbl.F_v,err = reader.ReadS32()
	checkErr(err)

	return
}

func PKT_MOVE(reader *packet.Packet)(tbl MOVE, err error){
	tbl.F_OID,err = reader.ReadS32()
	checkErr(err)

	tbl.F_X,err = reader.ReadS32()
	checkErr(err)

	tbl.F_Y,err = reader.ReadS32()
	checkErr(err)

	return
}

