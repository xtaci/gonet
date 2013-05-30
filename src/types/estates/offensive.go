package estates

import (
	"fmt"
	"time"
)

type Offensive struct {
	TYPE   string // Object Type
	OID    uint32 // Object ID
	X      uint16 // coordinate X
	Y      uint16 // coordinate Y
	Level  uint8
	Status byte
	Spec   map[string]string
}

type OffManager struct {
	Offensives []*Offensive
	CDs        map[string]*CD
}

func (m *OffManager) Append(estate *Offensive) {
	m.Offensives = append(m.Offensives, estate)
}

func (m *OffManager) AppendCD(event_id uint32, cd *CD) {
	if m.CDs == nil {
		m.CDs = make(map[string]*CD)
	}
	m.CDs[fmt.Sprint(event_id)] = cd
}

//------------------------------------------------ return num of changes
func (m *OffManager) CheckCD() int {
	opcount := 0
	for i := range m.CDs {
		if m.CDs[i].Timeout <= time.Now().Unix() { // times up
			for k := range m.Offensives {
				if m.CDs[i].OID == m.Offensives[k].OID { // if it is the oid
					m.Offensives[k].Status = STATUS_NORMAL
					opcount++
				}
			}
			delete(m.CDs, i)
		}
	}

	return opcount
}
