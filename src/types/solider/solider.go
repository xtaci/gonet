package solider

type Solider struct {
	OID      uint32 // object id
	Property interface{}
}

type SoliderCD struct {
	OID     uint32
	EventId uint32
}
