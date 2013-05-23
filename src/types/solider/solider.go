package solider

type Solider struct {
	TYPE string // type in string
	OID  uint32 // object id
	HP   uint16
}

type SoliderCD struct {
	OID     uint32
	EventId uint32
}
