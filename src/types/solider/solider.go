package solider

import (
	"fmt"
)

type Solider struct {
	TYPE     string // type in string
	OID      uint32 // object id
	HP       int16
	Ready    bool
	Property map[string]string // unit's private data
}

type Manager struct {
	UserId   int32
	Version  uint32
	Soliders map[string]*Solider // OID->Solider
}

func (m *Manager) Append(oid uint32, solider *Solider) {
	if m.Soliders == nil {
		m.Soliders = make(map[string]*Solider)
	}

	m.Soliders[fmt.Sprint(oid)] = solider
}
