package types

//---------------------------------------------------------- 单条聊天记录
type Words struct {
	Words     string
	SpeakerId int32
	Speaker   string
	Timestamp int64
}
