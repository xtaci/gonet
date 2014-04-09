package ipc_service

import "misc/packet"

type user_login_info struct {
	F_mac_addr       string
	F_client_version int32
	F_new_user       bool
	F_user_name      string
}

type user_snapshot struct {
	F_id             int32
	F_name           string
	F_rank           int32
	F_archives       string
	F_protect_time   int32
	F_last_save_time int32
	F_server_time    int32
}

type command_result_pack struct {
	F_rst int32
}

type talk struct {
	F_user string
	F_msg  string
}

func PKT_user_login_info(reader *packet.Packet) (tbl user_login_info, err error) {
	tbl.F_mac_addr, err = reader.ReadString()
	checkErr(err)

	tbl.F_client_version, err = reader.ReadS32()
	checkErr(err)

	tbl.F_new_user, err = reader.ReadBool()
	checkErr(err)

	tbl.F_user_name, err = reader.ReadString()
	checkErr(err)

	return
}

func PKT_user_snapshot(reader *packet.Packet) (tbl user_snapshot, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)

	tbl.F_name, err = reader.ReadString()
	checkErr(err)

	tbl.F_rank, err = reader.ReadS32()
	checkErr(err)

	tbl.F_archives, err = reader.ReadString()
	checkErr(err)

	tbl.F_protect_time, err = reader.ReadS32()
	checkErr(err)

	tbl.F_last_save_time, err = reader.ReadS32()
	checkErr(err)

	tbl.F_server_time, err = reader.ReadS32()
	checkErr(err)

	return
}

func PKT_command_result_pack(reader *packet.Packet) (tbl command_result_pack, err error) {
	tbl.F_rst, err = reader.ReadS32()
	checkErr(err)

	return
}

func PKT_talk(reader *packet.Packet) (tbl talk, err error) {
	tbl.F_user, err = reader.ReadString()
	checkErr(err)

	tbl.F_msg, err = reader.ReadString()
	checkErr(err)

	return
}
