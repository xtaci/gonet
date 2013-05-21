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
	SX       uint8       // size in X coord
	SY       uint8       // size in Y coord
	Property interface{} // related property
}

type EstateCD struct {
	OID     uint32
	EventId uint32
}

type EstateManager struct {
	Estates []Estate
	CDs     []EstateCD
}
