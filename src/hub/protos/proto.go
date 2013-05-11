package protos

import "misc/packet"

type msg struct {
	F_id int32
	F_data []byte
}

type id struct {
	F_id int32
}

type changescore struct {
	F_id int32
	F_oldscore int32
	F_newscore int32
}

type getlist struct {
	F_A int32
	F_B int32
}

type id_score struct {
	F_id int32
	F_score int32
}

type getlist_result struct {
	F_items []id_score
}

type longresult struct {
	F_v int64
}

type stringresult struct {
	F_v string
}

type intresult struct {
	F_v int32
}

func pktread_msg(reader *packet.Packet)(tbl msg, err error){
	tbl.F_id,err = reader.ReadS32()
	checkErr(err)
	tbl.F_data,err = reader.ReadBytes()
	checkErr(err)
	return
}

func pktread_id(reader *packet.Packet)(tbl id, err error){
	tbl.F_id,err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_changescore(reader *packet.Packet)(tbl changescore, err error){
	tbl.F_id,err = reader.ReadS32()
	checkErr(err)
	tbl.F_oldscore,err = reader.ReadS32()
	checkErr(err)
	tbl.F_newscore,err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_getlist(reader *packet.Packet)(tbl getlist, err error){
	tbl.F_A,err = reader.ReadS32()
	checkErr(err)
	tbl.F_B,err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_id_score(reader *packet.Packet)(tbl id_score, err error){
	tbl.F_id,err = reader.ReadS32()
	checkErr(err)
	tbl.F_score,err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_getlist_result(reader *packet.Packet)(tbl getlist_result, err error){
	narr,err2 := reader.ReadU16()
	checkErr(err2)
	tbl.F_items=make([]id_score,narr)
	for i:=0;i<int(narr);i++ {
		tbl.F_items[i], err = pktread_id_score(reader)
	}
	return
}

func pktread_longresult(reader *packet.Packet)(tbl longresult, err error){
	return
}

func pktread_stringresult(reader *packet.Packet)(tbl stringresult, err error){
	tbl.F_v,err = reader.ReadString()
	checkErr(err)
	return
}

func pktread_intresult(reader *packet.Packet)(tbl intresult, err error){
	tbl.F_v,err = reader.ReadS32()
	checkErr(err)
	return
}

