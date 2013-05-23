package estate

import (
	"encoding/json"
	"sync/atomic"
)

const (
	STATUS_NORMAL     = 0
	STATUS_UPGRADING  = 1
	STATUS_RECRUITING = 2
)

type Estate struct {
	TYPE    string // Object Type
	OID     uint32 // Object ID
	X       uint16 // coordinate X
	Y       uint16 // coordinate Y
	Level   uint8
	Status  uint8
}

//----------------------------------------------- Estate Move event records
type Move struct {
	OID uint32
	X   uint16
	Y   uint16
}

//----------------------------------------------- Estate Cooldown event records
type CD struct {
	OID     uint32
	EventId uint32
}

type Manager struct {
	Estates []Estate
	CDs     []CD
	NextVal int32
}

func (m *Manager) JSON() string {
	val, _ := json.Marshal(m)
	return string(val)
}

func (m *Manager) GENID() int32 {
	return atomic.AddInt32(&m.NextVal, 1)
}
