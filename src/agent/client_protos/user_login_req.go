package protos

import (
	//	"strconv"
	"time"

//	"log"
)

import (
	//	"agent/ipc"
	//	"cfg"
	//	"db/user_tbl"
	"misc/packet"
	. "types"
)

var EPOCH = time.Unix(0, 0)

func P_user_login_req(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	/*
		tbl, _ := PKT_user_login_info(reader)
		writer := packet.Writer()
		failed := command_result_pack{F_rst: 0}
		success := user_snapshot{}
		//------------------------------------------------

		config := cfg.Get()
		version, _ := strconv.Atoi(config["version"])

		if tbl.F_client_version != int32(version) {
			ret = packet.Pack(Code["user_login_faild_ack"], failed, writer)
			return
		}

		// if new user, register to db & online user
		// or else do a login procedure
		if tbl.F_new_user == 0 {
			if user_tbl.LoginMAC(sess.User.Mac, &sess.User) {
				ipc.RegisterOnline(sess, sess.User.Id)
				_fill_user_snapshot(&sess.User, &success)
				ret = packet.Pack(Code["user_login_succeed_ack"], success, writer)
				return
			} else {
				ret = packet.Pack(Code["user_login_faild_ack"], failed, writer)
				return
			}
		} else {
			sess.User.Name = tbl.F_user_name
			sess.User.Mac = tbl.F_mac_addr
			sess.User.CreatedAt = time.Now()

			if user_tbl.New(&sess.User) {
				if ipc.AddUser(sess.User.Id) {
					_fill_user_snapshot(&sess.User, &success)
					ret = packet.Pack(Code["user_login_succeed_ack"], success, writer)
					return
				} else {
					ret = packet.Pack(Code["user_login_faild_ack"], failed, writer)
					log.Println("ipc AddUser failed!!!\n")
					return
				}
			} else {
				ret = packet.Pack(Code["user_login_faild_ack"], failed, writer)
				return
			}
		}

	*/
	return
}

/*
func _fill_user_snapshot(basic *Basic, snapshot *user_snapshot) {
	snapshot.F_id = basic.Id
	snapshot.F_name = basic.Name
	snapshot.F_rank = basic.Score

	info, _:= ipc.GetInfo(basic.Id)
	pt := info.ProtectTime - time.Now().Unix()
	if pt > 0 {
		snapshot.F_protect_time = int32(pt)
	} else {
		snapshot.F_protect_time = 0
	}

	snapshot.F_last_save_time = int32(user.LastSaveTime.Unix())
	snapshot.F_server_time = int32(time.Now().Unix())
}
*/
