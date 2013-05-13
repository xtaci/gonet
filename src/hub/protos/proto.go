package protos

import "misc/packet"

type MSG struct {
	F_id int32
	F_data []byte
}

type ID struct {
	F_id int32
}

type CHGSCORE struct {
	F_id int32
	F_oldscore int32
	F_newscore int32
}

type GETLIST struct {
	F_A int32
	F_B int32
}

type ID_SCORE struct {
	F_id int32
	F_score int32
}

type LIST struct {
	F_items []ID_SCORE
}

type LONG struct {
	F_v int64
}

type STRING struct {
	F_v string
}

type INT struct {
	F_v int32
}

func PKT_MSG(reader *packet.Packet)(tbl MSG, err error){
	tbl.F_id,err = reader.ReadS32()
	checkErr(err)
	tbl.F_data,err = reader.ReadBytes()
	checkErr(err)
	return
}

func PKT_ID(reader *packet.Packet)(tbl ID, err error){
	tbl.F_id,err = reader.ReadS32()
	checkErr(err)
	return
}

func PKT_CHGSCORE(reader *packet.Packet)(tbl CHGSCORE, err error){
	tbl.F_id,err = reader.ReadS32()
	checkErr(err)
	tbl.F_oldscore,err = reader.ReadS32()
	checkErr(err)
	tbl.F_newscore,err = reader.ReadS32()
	checkErr(err)
	return
}

func PKT_GETLIST(reader *packet.Packet)(tbl GETLIST, err error){
	tbl.F_A,err = reader.ReadS32()
	checkErr(err)
	tbl.F_B,err = reader.ReadS32()
	checkErr(err)
	return
}

func PKT_ID_SCORE(reader *packet.Packet)(tbl ID_SCORE, err error){
	tbl.F_id,err = reader.ReadS32()
	checkErr(err)
	tbl.F_score,err = reader.ReadS32()
	checkErr(err)
	return
}

func PKT_LIST(reader *packet.Packet)(tbl LIST, err error){
	narr,err2 := reader.ReadU16()
	checkErr(err2)
	tbl.F_items=make([]ID_SCORE,narr)
	for i:=0;i<int(narr);i++ {
		tbl.F_items[i], err = PKT_ID_SCORE(reader)
	}
	return
}

func PKT_LONG(reader *packet.Packet)(tbl LONG, err error){
	return
}

func PKT_STRING(reader *packet.Packet)(tbl STRING, err error){
	tbl.F_v,err = reader.ReadString()
	checkErr(err)
	return
}

func PKT_INT(reader *packet.Packet)(tbl INT, err error){
	tbl.F_v,err = reader.ReadS32()
	checkErr(err)
	return
}

