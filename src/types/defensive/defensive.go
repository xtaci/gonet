package defensive

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	STATUS_NORMAL = 0
	STATUS_CD     = 1
)

type Defensive struct {
	TYPE   string // Object Type
	OID    uint32 // Object ID
	X      uint16 // coordinate X
	Y      uint16 // coordinate Y
	Level  uint8
	Status uint8
}

//----------------------------------------------- Defensive Move event records
type Move struct {
	OID uint32
	X   uint16
	Y   uint16
}

//----------------------------------------------- Defensive Cooldown event records
type CD struct {
	OID     uint32
	Timeout int64
}

type Manager struct {
	Id      int32
	Defensives []Defensive
	CDs     map[string]*CD
	NextVal uint32
	Version uint32
}

func (m *Manager) AppendDefensive(estate *Defensive) {
	m.Defensives = append(m.Defensives, *estate)
}

func (m *Manager) AppendCD(event_id uint32, cd *CD) {
	if m.CDs == nil {
		m.CDs = make(map[string]*CD)
	}
	m.CDs[fmt.Sprint(event_id)] = cd
}

func (m *Manager) GENID() uint32 {
	return atomic.AddUint32(&m.NextVal, 1)
}

func (m *Manager) CheckCD() int {
	opcount := 0
	for i := range m.CDs {
		if m.CDs[i].Timeout <= time.Now().Unix() { // times up
			for k := range m.Defensives {
				if m.CDs[i].OID == m.Defensives[k].OID { // if it is the oid
					m.Defensives[k].Status = STATUS_NORMAL
					opcount++
				}
			}
			delete(m.CDs, i)
		}
	}

	return opcount
}
