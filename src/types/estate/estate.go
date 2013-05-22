package estate

const (
	TYPE_DEFENSE  = uint16(1)
	TYPE_BARRACKS = uint16(2)
)

type Estate struct {
	OID      uint32      // Object ID
	TYPE     uint16      // Object Type
	X        uint16      // coordinate X
	Y        uint16      // coordinate Y
	Property interface{} // related property
}

//----------------------------------------------- Estate Move event records
type Move struct {
	OID uint32
	X	uint16
	Y	uint16
}

//----------------------------------------------- Estate Cooldown event records
type CD struct {
	OID     uint32
	EventId uint32
}

type Manager struct {
	Estates []Estate
	CDs     []CD
}
