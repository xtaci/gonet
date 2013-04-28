package protos

import (
	"strconv"
)

import (
	. "types"
	"misc/packet"
	"hub/names"
	"hub/ranklist"
	"db/user_tbl"
	"cfg"
)

func _user_login_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {

	tbl,_ := pktread_user_login_info(reader)
	writer := packet.Writer()
	failed := command_result_pack{rst:0}

	config := cfg.Get()
	version, _ :=  strconv.Atoi(config["version"])

	if tbl.client_version != int32(version) {
		ret = pack(failed, writer)
		return
	}

	if tbl.new_user == 0 {
		if user_tbl.LoginMAC(sess.User.Mac, &sess.User) {
			names.Register(sess.MQ, sess.User.Id)
			success := user_snapshot{}
			success.id = sess.User.Id
			success.name = sess.User.Name
			success.rank = sess.User.Score
			ret = pack(success, writer)
			return
		} else {
			ret = pack(failed,writer)
			return
		}
	} else {
		sess.User.Name = tbl.user_name
		sess.User.Mac = tbl.mac_addr
		sess.User.Score = ranklist.Increase()

		if user_tbl.New(&sess.User) {
			names.Register(sess.MQ, sess.User.Id)
			success := user_snapshot{}
			ret = pack(success,writer)
			return
		} else {
			ranklist.Decrease()
			ret = pack(failed,writer)
			return
		}
	}

	return nil,nil
}
