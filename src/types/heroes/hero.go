package hero

import (
	"fmt"
	"time"
)

const (
	COLLECTION = "HEROES"
	STATUS_NORMAL = byte(0)
	STATUS_CD     = 1
)

//----------------------------------------------- Generic Cooldown event records
type CD struct {
	OID     uint32
	Timeout int64
}

type Hero struct {
	TYPE string // type in string
	Status byte
	Spec map[string]string
}

type Manager struct {
	UserId  int32
	Version uint32
	Heroes  map[string]*Hero // OID->Hero
	CDs     map[string]*CD   // waiting CoolDown
}

func (m *Manager) Append(oid uint32, hero *Hero) {
	if m.Heroes == nil {
		m.Heroes = make(map[string]*Hero)
	}

	m.Heroes[fmt.Sprint(oid)] = hero
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
			if hero := m.Heroes[oid]; hero != nil {
				hero.Status = STATUS_NORMAL
				opcount++
				delete(m.Heroes, oid)
			}
			delete(m.CDs, k)
		}
	}

	return opcount
}
