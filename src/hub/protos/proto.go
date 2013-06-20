package protos

import "misc/packet"

type FORWARDIPC struct {
	F_IPC []byte
}

type LOGIN_REQ struct {
	F_id          int32
	F_group       int32
	F_groupmsgmax uint32
}

type LOGIN_ACK struct {
	F_success     bool
	F_groupmsgmax uint32
}

type ID struct {
	F_id int32
}

type PROTECT struct {
	F_id          int32
	F_protecttime int64
}

type INFO struct {
	F_flag        bool
	F_id          int32
	F_state       byte
	F_score       int32
	F_protecttime int64
}

type CHGSCORE struct {
	F_id       int32
	F_oldscore int32
	F_newscore int32
}

type GETLIST struct {
	F_A int32
	F_B int32
}

type ID_SCORE struct {
	F_id    int32
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

func PKT_FORWARDIPC(reader *packet.Packet) (tbl FORWARDIPC, err error) {
	tbl.F_IPC, err = reader.ReadBytes()
	checkErr(err)

	return
}

func PKT_LOGIN_REQ(reader *packet.Packet) (tbl LOGIN_REQ, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)

	tbl.F_group, err = reader.ReadS32()
	checkErr(err)

	tbl.F_groupmsgmax, err = reader.ReadU32()
	checkErr(err)

	return
}

func PKT_LOGIN_ACK(reader *packet.Packet) (tbl LOGIN_ACK, err error) {
	tbl.F_success, err = reader.ReadBool()
	checkErr(err)

	tbl.F_groupmsgmax, err = reader.ReadU32()
	checkErr(err)

	return
}

func PKT_ID(reader *packet.Packet) (tbl ID, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)

	return
}

func PKT_PROTECT(reader *packet.Packet) (tbl PROTECT, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)

	tbl.F_protecttime, err = reader.ReadS64()
	checkErr(err)

	return
}

func PKT_INFO(reader *packet.Packet) (tbl INFO, err error) {
	tbl.F_flag, err = reader.ReadBool()
	checkErr(err)

	tbl.F_id, err = reader.ReadS32()
	checkErr(err)

	tbl.F_state, err = reader.ReadByte()
	checkErr(err)

	tbl.F_score, err = reader.ReadS32()
	checkErr(err)

	tbl.F_protecttime, err = reader.ReadS64()
	checkErr(err)

	return
}

func PKT_CHGSCORE(reader *packet.Packet) (tbl CHGSCORE, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)

	tbl.F_oldscore, err = reader.ReadS32()
	checkErr(err)

	tbl.F_newscore, err = reader.ReadS32()
	checkErr(err)

	return
}

func PKT_GETLIST(reader *packet.Packet) (tbl GETLIST, err error) {
	tbl.F_A, err = reader.ReadS32()
	checkErr(err)

	tbl.F_B, err = reader.ReadS32()
	checkErr(err)

	return
}

func PKT_ID_SCORE(reader *packet.Packet) (tbl ID_SCORE, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)

	tbl.F_score, err = reader.ReadS32()
	checkErr(err)

	return
}

func PKT_LIST(reader *packet.Packet) (tbl LIST, err error) {
	narr := uint16(0)

	narr, err = reader.ReadU16()
	checkErr(err)

	tbl.F_items = make([]ID_SCORE, narr)
	for i := 0; i < int(narr); i++ {
		tbl.F_items[i], err = PKT_ID_SCORE(reader)
		checkErr(err)

	}

	return
}

func PKT_LONG(reader *packet.Packet) (tbl LONG, err error) {
	tbl.F_v, err = reader.ReadS64()
	checkErr(err)

	return
}

func PKT_STRING(reader *packet.Packet) (tbl STRING, err error) {
	tbl.F_v, err = reader.ReadString()
	checkErr(err)

	return
}

func PKT_INT(reader *packet.Packet) (tbl INT, err error) {
	tbl.F_v, err = reader.ReadS32()
	checkErr(err)

	return
}
