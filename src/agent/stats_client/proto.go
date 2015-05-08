package stats_client

import "misc/packet"

type SET_ADDS_REQ struct {
	F_key   string
	F_value int32
	F_lang  string
}

type SET_UPDATE_REQ struct {
	F_key   string
	F_value string
	F_lang  string
}

func PKT_SET_ADDS_REQ(reader *packet.Packet) (tbl SET_ADDS_REQ, err error) {
	tbl.F_key, err = reader.ReadString()
	checkErr(err)

	tbl.F_value, err = reader.ReadS32()
	checkErr(err)

	tbl.F_lang, err = reader.ReadString()
	checkErr(err)

	return
}

func PKT_SET_UPDATE_REQ(reader *packet.Packet) (tbl SET_UPDATE_REQ, err error) {
	tbl.F_key, err = reader.ReadString()
	checkErr(err)

	tbl.F_value, err = reader.ReadString()
	checkErr(err)

	tbl.F_lang, err = reader.ReadString()
	checkErr(err)

	return
}
