package protos

import "misc/packet"

type user_login_info struct {
	mac_addr string
	client_version int32
	new_user byte
	user_name string
}

type user_snapshot struct {
	id int32
	name string
	rank int32
	archives string
	protect_time int32
	last_save_time int32
	server_time int32
}

type command_result_pack struct {
	rst int32
}

type user_archives_info struct {
	id int32
	archives string
}

type rank_list_item struct {
	id int32
	name string
	rank int32
	state int32
	protect_time int32
}

type rank_list struct {
	items []rank_list_item
}

type pve_list_item struct {
	id int32
	name string
	rank int32
	state int32
	protect_time int32
}

type pve_list struct {
	items []pve_list_item
}

type command_id_pack struct {
	id int32
}

type atk_player_rst_req struct {
	rst int32
	protect_time int32
}

type atk_monster_rst_req struct {
	protect_time int32
}

func pktread_user_login_info(reader *packet.Packet)(tbl user_login_info, err error){
	tbl.mac_addr,err = reader.ReadString()
	tbl.client_version,err = reader.ReadS32()
	checkErr(err)
	tbl.new_user,err = reader.ReadByte()
	checkErr(err)
	tbl.user_name,err = reader.ReadString()
	return
}

func pktread_user_snapshot(reader *packet.Packet)(tbl user_snapshot, err error){
	tbl.id,err = reader.ReadS32()
	checkErr(err)
	tbl.name,err = reader.ReadString()
	tbl.rank,err = reader.ReadS32()
	checkErr(err)
	tbl.archives,err = reader.ReadString()
	tbl.protect_time,err = reader.ReadS32()
	checkErr(err)
	tbl.last_save_time,err = reader.ReadS32()
	checkErr(err)
	tbl.server_time,err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_command_result_pack(reader *packet.Packet)(tbl command_result_pack, err error){
	tbl.rst,err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_user_archives_info(reader *packet.Packet)(tbl user_archives_info, err error){
	tbl.id,err = reader.ReadS32()
	checkErr(err)
	tbl.archives,err = reader.ReadString()
	return
}

func pktread_rank_list_item(reader *packet.Packet)(tbl rank_list_item, err error){
	tbl.id,err = reader.ReadS32()
	checkErr(err)
	tbl.name,err = reader.ReadString()
	tbl.rank,err = reader.ReadS32()
	checkErr(err)
	tbl.state,err = reader.ReadS32()
	checkErr(err)
	tbl.protect_time,err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_rank_list(reader *packet.Packet)(tbl rank_list, err error){
	narr,err2 := reader.ReadU16()
	checkErr(err2)
	tbl.items=make([]rank_list_item,narr)
	for i:=0;i<int(narr);i++ {
		tbl.items[i], err = pktread_rank_list_item(reader)
	}
	return
}

func pktread_pve_list_item(reader *packet.Packet)(tbl pve_list_item, err error){
	tbl.id,err = reader.ReadS32()
	checkErr(err)
	tbl.name,err = reader.ReadString()
	tbl.rank,err = reader.ReadS32()
	checkErr(err)
	tbl.state,err = reader.ReadS32()
	checkErr(err)
	tbl.protect_time,err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_pve_list(reader *packet.Packet)(tbl pve_list, err error){
	narr,err2 := reader.ReadU16()
	checkErr(err2)
	tbl.items=make([]pve_list_item,narr)
	for i:=0;i<int(narr);i++ {
		tbl.items[i], err = pktread_pve_list_item(reader)
	}
	return
}

func pktread_command_id_pack(reader *packet.Packet)(tbl command_id_pack, err error){
	tbl.id,err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_atk_player_rst_req(reader *packet.Packet)(tbl atk_player_rst_req, err error){
	tbl.rst,err = reader.ReadS32()
	checkErr(err)
	tbl.protect_time,err = reader.ReadS32()
	checkErr(err)
	return
}

func pktread_atk_monster_rst_req(reader *packet.Packet)(tbl atk_monster_rst_req, err error){
	tbl.protect_time,err = reader.ReadS32()
	checkErr(err)
	return
}

