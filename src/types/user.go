package types

type User struct {
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
