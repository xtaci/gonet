package estates

import (
	"fmt"
	"sync/atomic"
)

const (
	COLLECTION = "ESTATES"
)

//----------------------------------------------- Generic Move event records
type Move struct {
	OID uint16
	X   uint8
	Y   uint8
}

type Estate struct {
	TYPE     uint32 // Object Type
	X        uint8  // coordinate X
	Y        uint8  // coordinate Y
	Ready    bool
	Property map[string]string // unit's private data
}

type Manager struct {
	UserId  int32
	Version uint32
	Estates map[string]*Estate // OID->Estate
	NextVal uint32
}

//------------------------------------------------ Generate Estate object-id
func (m *Manager) NextID() uint16 {
	v := atomic.AddUint32(&m.NextVal, 1)
	oid := uint16(v)

	// collision? loop!
	if m.Estates[fmt.Sprint(oid)] != nil {
		return m.NextID()
	}

	return oid
}

func (m *Manager) Append(oid uint16, estate *Estate) {
	if m.Estates == nil {
		m.Estates = make(map[string]*Estate)
	}

	m.Estates[fmt.Sprint(oid)] = estate
}
