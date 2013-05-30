package estates

import (
	"fmt"
	"time"
)

type Basic struct {
	TYPE   string // Object Type
	OID    uint32 // Object ID
	X      uint16 // coordinate X
	Y      uint16 // coordinate Y
	Level  uint8
	Status byte
}

type BasManager struct {
	Basics []*Basic
	CDs    map[string]*CD
}

func (m *BasManager) Append(estate *Basic) {
	m.Basics = append(m.Basics, estate)
}

func (m *BasManager) AppendCD(event_id uint32, cd *CD) {
	if m.CDs == nil {
		m.CDs = make(map[string]*CD)
	}
	m.CDs[fmt.Sprint(event_id)] = cd
}

//------------------------------------------------ return num of changes
func (m *BasManager) CheckCD() int {
	opcount := 0
	for i := range m.CDs {
		if m.CDs[i].Timeout <= time.Now().Unix() { // times up
			for k := range m.Basics {
				if m.CDs[i].OID == m.Basics[k].OID { // if it is the oid
					m.Basics[k].Status = STATUS_NORMAL
					opcount++
				}
			}
			delete(m.CDs, i)
		}
	}

	return opcount
}
