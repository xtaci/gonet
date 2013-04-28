package protos

import "misc/packet"

type user_login_info struct {
	F_mac_addr       string
	F_client_version int32
	F_new_user       byte
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

type user_archives_info struct {
	F_id       int32
	F_archives string
}

type rank_list_item struct {
	F_id           int32
	F_name         string
	F_rank         int32
	F_state        int32
	F_protect_time int32
}

type rank_list struct {
	F_items []rank_list_item
}

type pve_list_item struct {
	F_id           int32
	F_name         string
	F_rank         int32
	F_state        int32
	F_protect_time int32
}

type pve_list struct {
	F_items []pve_list_item
}

type command_id_pack struct {
	F_id int32
}

type atk_player_rst_req struct {
	F_rst          int32
	F_protect_time int32
}

type atk_monster_rst_req struct {
	F_protect_time int32
}

func pktread_user_login_info(reader *packet.Packet) (tbl user_login_info, err error) {
	tbl.F_mac_addr, err = reader.ReadString()
	tbl.F_client_version, err = reader.ReadS32()
	checkErr(err)
	tbl.F_new_user, err = reader.ReadByte()
	checkErr(err)
	tbl.F_user_name, err = reader.ReadString()
	return
}

func pktread_user_snapshot(reader *packet.Packet) (tbl user_snapshot, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)
	tbl.F_name, err = reader.ReadString()
	tbl.F_rank, err = reader.ReadS32()
	checkErr(err)
	tbl.F_archives, err = reader.ReadString()
	tbl.F_protect_time, err = reader.ReadS32()
	checkErr(err)
	tbl.F_last_save_time, err = reader.ReadS32()
	checkErr(err)
	tbl.F_server_time, err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_command_result_pack(reader *packet.Packet) (tbl command_result_pack, err error) {
	tbl.F_rst, err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_user_archives_info(reader *packet.Packet) (tbl user_archives_info, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)
	tbl.F_archives, err = reader.ReadString()
	return
}

func pktread_rank_list_item(reader *packet.Packet) (tbl rank_list_item, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)
	tbl.F_name, err = reader.ReadString()
	tbl.F_rank, err = reader.ReadS32()
	checkErr(err)
	tbl.F_state, err = reader.ReadS32()
	checkErr(err)
	tbl.F_protect_time, err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_rank_list(reader *packet.Packet) (tbl rank_list, err error) {
	narr, err2 := reader.ReadU16()
	checkErr(err2)
	tbl.F_items = make([]rank_list_item, narr)
	for i := 0; i < int(narr); i++ {
		tbl.F_items[i], err = pktread_rank_list_item(reader)
	}
	return
}

func pktread_pve_list_item(reader *packet.Packet) (tbl pve_list_item, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)
	tbl.F_name, err = reader.ReadString()
	tbl.F_rank, err = reader.ReadS32()
	checkErr(err)
	tbl.F_state, err = reader.ReadS32()
	checkErr(err)
	tbl.F_protect_time, err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_pve_list(reader *packet.Packet) (tbl pve_list, err error) {
	narr, err2 := reader.ReadU16()
	checkErr(err2)
	tbl.F_items = make([]pve_list_item, narr)
	for i := 0; i < int(narr); i++ {
		tbl.F_items[i], err = pktread_pve_list_item(reader)
	}
	return
}

func pktread_command_id_pack(reader *packet.Packet) (tbl command_id_pack, err error) {
	tbl.F_id, err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_atk_player_rst_req(reader *packet.Packet) (tbl atk_player_rst_req, err error) {
	tbl.F_rst, err = reader.ReadS32()
	checkErr(err)
	tbl.F_protect_time, err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_atk_monster_rst_req(reader *packet.Packet) (tbl atk_monster_rst_req, err error) {
	tbl.F_protect_time, err = reader.ReadS32()
	checkErr(err)
	return
}
