package types

type User struct {
	Id             int32  // user id
	ClanId         int32  // clan id
	ClanMsgId      uint32 // clan's message queue sequence id
	Name           string // player name
	Pass           []byte // player password
	Mac            string // player's mac address
	Score          int32  // player's score
	ProtectTimeout int64
	LoginCount     int32
	LastLogin      int64
}
