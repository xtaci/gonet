package estates

import (
	"fmt"
	"time"
)

type Defensive struct {
	TYPE   string // Object Type
	OID    uint32 // Object ID
	X      uint16 // coordinate X
	Y      uint16 // coordinate Y
	Level  uint8
	Status byte
}

type DefManager struct {
	Defensives []*Defensive
	CDs        map[string]*CD
}

func (m *DefManager) Append(estate *Defensive) {
	m.Defensives = append(m.Defensives, estate)
}

func (m *DefManager) AppendCD(event_id uint32, cd *CD) {
	if m.CDs == nil {
		m.CDs = make(map[string]*CD)
	}
	m.CDs[fmt.Sprint(event_id)] = cd
}

//------------------------------------------------ return num of changes
func (m *DefManager) CheckCD() int {
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
