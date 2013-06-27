package ipc

import (
	"math/big"
)

import (
	"agent/AI"
	"db/user_tbl"
	"misc/crypto/diffie"
	"misc/crypto/pike"
	"misc/packet"
	. "types"
)

func P_heart_beat_req(sess *Session, reader *packet.Packet) []byte {
	// nothing should be done
	return nil
}

func P_user_login_req(sess *Session, reader *packet.Packet) []byte {
	tbl, _ := PKT_user_login_info(reader)
	ret := command_result_pack{}

	if user := user_tbl.LoginMac(tbl.F_user_name, tbl.F_mac_addr); user != nil {
		sess.User = user
		AI.LoginProc(sess)
	}

	return packet.Pack(Code["user_login_succeed_ack"], &ret, nil)
}

func P_talk_req(sess *Session, reader *packet.Packet) []byte {
	tbl, _ := PKT_talk(reader)
	dest := user_tbl.Query(tbl.F_user)
	Send(sess.User.Id, dest.Id, SERVICE_TALK, false, tbl.F_msg)
	return nil
}

func P_key_exchange_req(sess *Session, reader *packet.Packet) []byte {
	tbl, _ := PKT_KEY(reader)
	A := big.NewInt(int64(tbl.F_E))
	secret, B := diffie.DHGenKey(diffie.DH1BASE, diffie.DH1PRIME)

	key := big.NewInt(0).Exp(A, secret, diffie.DH1PRIME)
	sess.Crypto = pike.NewCtx(uint32(key.Uint64()))

	ret := KEY{int32(B.Uint64())}
	return packet.Pack(Code["key_exchange_ack"], &ret, nil)
}
