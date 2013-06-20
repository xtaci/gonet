package types

type Event struct {
	UserId  int32
	EventId int32
	Type    int16
	Params  []byte
	Timeout int64
}
