package hero

import (
	"fmt"
)

const (
	COLLECTION = "HEROES"
)

type Hero struct {
	TYPE     string // type in string
	HP       int16
	Ready    bool
	Property map[string]string // unit's private data
}

type Manager struct {
	UserId  int32
	Version uint32
	Heroes  map[string]*Hero // OID->Hero
}

func (m *Manager) Append(oid uint32, hero *Hero) {
	if m.Heroes == nil {
		m.Heroes = make(map[string]*Hero)
	}

	m.Heroes[fmt.Sprint(oid)] = hero
}
