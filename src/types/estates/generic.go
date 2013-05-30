package estates

const (
	STATUS_NORMAL = byte(0)
	STATUS_CD     = 1

	CDTYPE_BASIC     = byte(1)
	CDTYPE_OFFENSIVE = 2
	CDTYPE_DEFENSIVE = 3
	CDTYPE_RESOURCE  = 4
)

//----------------------------------------------- Generic Move event records
type Move struct {
	OID uint32
	X   uint16
	Y   uint16
}

//----------------------------------------------- Generic Cooldown event records
type CD struct {
	OID     uint32
	Timeout int64
}
