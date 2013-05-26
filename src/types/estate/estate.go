package estate

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

const (
	STATUS_NORMAL = 0
	STATUS_CD     = 1
)

type Estate struct {
	TYPE   string // Object Type
	OID    uint32 // Object ID
	X      uint16 // coordinate X
	Y      uint16 // coordinate Y
	Level  uint8
	Status uint8
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
	Timeout int64
}

type Manager struct {
	Id      int32
	Estates []Estate
	CDs     map[string]*CD
	NextVal uint32
	Version uint32
}

func (m *Manager) AppendEstate(estate *Estate) {
	m.Estates = append(m.Estates, *estate)
}

func (m *Manager) AppendCD(event_id uint32, cd *CD) {
	if m.CDs == nil {
		m.CDs = make(map[string]*CD)
	}
	m.CDs[fmt.Sprint(event_id)] = cd
}

func (m *Manager) JSON() string {
	val, _ := json.Marshal(m)
	return string(val)
}

func (m *Manager) GENID() uint32 {
	return atomic.AddUint32(&m.NextVal, 1)
}
