package estates

import (
	"fmt"
	"time"
)

type Resource struct {
	TYPE   string // Object Type
	OID    uint32 // Object ID
	X      uint16 // coordinate X
	Y      uint16 // coordinate Y
	Level  uint8
	Status byte
	Spec   map[string]string
}

type ResManager struct {
	Resources []*Resource
	CDs       map[string]*CD
}

func (m *ResManager) Append(estate *Resource) {
	m.Resources = append(m.Resources, estate)
}

func (m *ResManager) AppendCD(event_id uint32, cd *CD) {
	if m.CDs == nil {
		m.CDs = make(map[string]*CD)
	}
	m.CDs[fmt.Sprint(event_id)] = cd
}

//------------------------------------------------ return num of changes
func (m *ResManager) CheckCD() int {
	opcount := 0
	for i := range m.CDs {
		if m.CDs[i].Timeout <= time.Now().Unix() { // times up
			for k := range m.Resources {
				if m.CDs[i].OID == m.Resources[k].OID { // if it is the oid
					m.Resources[k].Status = STATUS_NORMAL
					opcount++
				}
			}
			delete(m.CDs, i)
		}
	}

	return opcount
}
