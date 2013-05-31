package estates

import (
	"fmt"
	"time"
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

//----------------------------------------------- Generic Cooldown event records
type CD struct {
	OID     uint32
	Timeout int64
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
	CDs     map[string]*CD     // EventId->CD
}

func (m *Manager) Append(oid uint32, estate *Estate) {
	if m.Estates == nil {
		m.Estates = make(map[string]*Estate)
	}

	m.Estates[fmt.Sprint(oid)] = estate
}

func (m *Manager) AppendCD(event_id uint32, cd *CD) {
	if m.CDs == nil {
		m.CDs = make(map[string]*CD)
	}
	m.CDs[fmt.Sprint(event_id)] = cd
}

//------------------------------------------------ return num of changes
func (m *Manager) CheckCD() int {
	opcount := 0
	for k := range m.CDs {
		if m.CDs[k].Timeout <= time.Now().Unix() { // times up
			oid := fmt.Sprint(m.CDs[k].OID)
			if estate := m.Estates[oid]; estate != nil {
				estate.Status = STATUS_NORMAL
				opcount++
				delete(m.Estates, oid)
			}
			delete(m.CDs, k)
		}
	}

	return opcount
}
