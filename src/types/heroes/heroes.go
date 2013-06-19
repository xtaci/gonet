package heroes

const (
	COLLECTION = "HEROES"
)

type Hero struct { // forcus on skills & equip
	Ready    bool
	Property map[string]string // hero's private data, like skills
}

type Manager struct {
	UserId  int32
	Version uint32
	Heroes  map[string]*Hero // Type->Hero
}

func (m *Manager) Unlock(Type string) {
	if m.Heroes == nil {
		m.Heroes = make(map[string]*Hero)
	}

	if m.Heroes[Type] == nil {
		m.Heroes[Type] = &Hero{Property: make(map[string]string)}
	}
}
