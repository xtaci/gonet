package soldiers

import (
	"fmt"
)

const (
	COLLECTION = "SOLDIERS"
)

type Soldier struct {
	TYPE     string // type in string
	OID      uint32 // object id
	HP       int16
	Ready    bool
	Property map[string]string // unit's private data
}

type Manager struct {
	UserId   int32
	Version  uint32
	Soldiers map[string]*Soldier // OID->Soldier
}

func (m *Manager) Append(oid uint32, soldier *Soldier) {
	if m.Soldiers == nil {
		m.Soldiers = make(map[string]*Soldier)
	}

	m.Soldiers[fmt.Sprint(oid)] = soldier
}
