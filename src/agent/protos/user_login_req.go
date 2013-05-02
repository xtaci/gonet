package protos

import (
	"strconv"
	"time"
)

import (
	"cfg"
	"db/user_tbl"
	"hub/names"
	"hub/ranklist"
	"misc/packet"
	. "types"
)

func _user_login_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {

	tbl, _ := pktread_user_login_info(reader)
	writer := packet.Writer()
	failed := command_result_pack{F_rst: 0}

	config := cfg.Get()
	version, _ := strconv.Atoi(config["version"])

	if tbl.F_client_version != int32(version) {
		ret = pack(Code["user_login_faild_ack"], failed, writer)
		return
	}

	if tbl.F_new_user == 0 {
		if user_tbl.LoginMAC(sess.User.Mac, &sess.User) {
			names.Register(sess, sess.User.Id)
			success := user_snapshot{}
			_fill_user_snapshot(&sess.User, &success)
			ret = pack(Code["user_login_succeed_ack"], success, writer)
			return
		} else {
			ret = pack(Code["user_login_faild_ack"], failed, writer)
			return
		}
	} else {
		sess.User.Name = tbl.F_user_name
		sess.User.Mac = tbl.F_mac_addr
		sess.User.Score = ranklist.Increase()
		sess.User.State = 0
		sess.User.LastSaveTime = time.Now()
		sess.User.ProtectTime = time.Now()
		sess.User.CreatedAt = time.Now()

		if user_tbl.New(&sess.User) {
			names.Register(sess, sess.User.Id)
			success := user_snapshot{}
			_fill_user_snapshot(&sess.User, &success)
			ret = pack(Code["user_login_succeed_ack"], success, writer)
			return
		} else {
			ranklist.Decrease()
			ret = pack(Code["user_login_faild_ack"], failed, writer)
			return
		}
	}

	return nil, nil
}

func _fill_user_snapshot(user *User, snapshot *user_snapshot) {
	snapshot.F_id = user.Id
	snapshot.F_name = user.Name
	snapshot.F_rank = user.Score

	pt := user.ProtectTime.Unix() - time.Now().Unix()
	if pt > 0 {
		snapshot.F_protect_time = int32(pt)
	} else {
		snapshot.F_protect_time = 0
	}

	snapshot.F_last_save_time = int32(user.LastSaveTime.Unix())
	snapshot.F_server_time = int32(time.Now().Unix())
	snapshot.F_archives = user.Archives
}
