package estates

import (
	"fmt"
)

const (
	COLLECTION    = "ESTATES"
	STATUS_NORMAL = byte(0)
	STATUS_CD     = 1
)

//----------------------------------------------- Generic Move event records
type Move struct {
	OID uint32
	X   uint8
	Y   uint8
}

type Estate struct {
	TYPE   uint32 // Object Type
	X      uint8  // coordinate X
	Y      uint8  // coordinate Y
	Status byte
	Spec   map[string]string // unit's private data
}

type Manager struct {
	UserId  int32
	Version uint32
	Estates map[string]*Estate // OID->Estate
}

func (m *Manager) Append(oid uint32, estate *Estate) {
	if m.Estates == nil {
		m.Estates = make(map[string]*Estate)
	}

	m.Estates[fmt.Sprint(oid)] = estate
}
