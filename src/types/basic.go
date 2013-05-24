package types

import (
	"encoding/json"
)

type Basic struct {
	Id             int32
	Name           string
	Pass           []byte
	Mac            string
	Score          int32
	ProtectTimeout int64
	IsProtecting   bool
	LoginCount     int32
	LastLogin      int64
}

func (b *Basic) JSON() string {
	val, _ := json.Marshal(b)
	return string(val)
}
