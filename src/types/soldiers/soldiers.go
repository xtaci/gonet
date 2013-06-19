package soldiers

const (
	COLLECTION = "SOLDIERS"
)

//------------------------------------------------ soldiers focus on quantity !!!
type Soldier struct {
	Count    int32             // num of soliders
	Property map[string]string // private data for this type of soldier
}

type Manager struct {
	UserId   int32
	Version  uint32
	Soldiers map[string]*Soldier // Type->Soldier
}

func (m *Manager) Add(Type string, num int32) {
	if m.Soldiers == nil {
		m.Soldiers = make(map[string]*Soldier)
	}

	if m.Soldiers[Type] == nil {
		m.Soldiers[Type] = &Soldier{Count: num, Property: make(map[string]string)}
	} else {
		m.Soldiers[Type].Count += num
	}
}

func (m *Manager) Remove(Type string, num int32) {
	if m.Soldiers != nil && m.Soldiers[Type] != nil {
		m.Soldiers[Type].Count -= num
	}
}
