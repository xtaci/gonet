package types

const (
	MAP_W = 40
	MAP_H = 40
)

type Map struct {
	Bitset []byte
}

func MapLoad(bitmap string) *Map {
	m := &Map{}
	m.Bitset = make([]byte, MAP_W*MAP_H)
	return m
}

func (m *Map) MapTest(X,Y int) bool {
	if X < MAP_W && Y < MAP_H {
		bit := Y * MAP_W + X
		n := bit/8
		off := uint(bit%8)
		if (m.Bitset[n] & (byte(128) >> off)) !=0 {
			 return true
		}
	}

	return false
}

func (m *Map) MapSet(X,Y int) bool {
	return false
}

func (m *Map) MapUnset(X,Y int) bool {
	return false
}
