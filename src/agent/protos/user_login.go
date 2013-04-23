package protos
/*
import . "types"
import "misc/packet"
import "hub/names"
import "db/user"
import "db/city"
import "time"
func UserLogin(sess *Session, reader *packet.Packet) (ret []byte, err error) {
	isLoggedIn := false
	if register == 0 {
		if user.LoginMAC(mac, &sess.User) {
			names.Register(sess.MQ, sess.User.Id)
			isLoggedIn = true
		}
	} else {
		name := reader.ReadString()
		checkErr(err)
		if user.New(name, mac, &sess.User) {
			names.Register(sess.MQ, sess.User.Id)
		}
	}

	writer := packet.PacketWriter()
	writer.WriteU16(2)
	writer.WriteU32(sess.User.Id)
	writer.WriteString(sess.User.Name)
	writer.WriteU32(sess.User.Score)
	writer.WriteString(sess.User.Data)		//玩家信息
	writer.WriteU32(sess.User.ShieldUtil.Unix() - time.Now().Unix())
	writer.WriteU32(sess.User.LastSync.Unix())
	writer.WriteU32(time.Now().Unix())

	return writer.Data(), err
}
*/
