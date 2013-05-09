package grid

const (
	GRID_W = 40
	GRID_H = 40
)

type Grid struct {
	Bitset []byte
}

//----------------------------------------------- Create a new grid struc
func NewGrid() *Grid {
	m := &Grid{}
	m.Bitset = make([]byte, GRID_W*GRID_H)
	return m
}

//----------------------------------------------- test whether a (X,Y) is set
func (m *Grid) Test(X,Y int) bool {
	if X < GRID_W && Y < GRID_H {
		bit := Y * GRID_W + X
		n := bit/8
		off := uint(bit%8)
		if (m.Bitset[n] & (byte(128) >> off)) !=0 {
			 return true
		}
	}

	return false
}

//----------------------------------------------- Set (X,Y) -> 1
func (m *Grid) Set(X,Y int) {
	if X < GRID_W && Y < GRID_H {
		bit := Y * GRID_W + X
		n := bit/8
		off := uint(bit%8)
		m.Bitset[n] |= byte(128)>>off
	}
}

//----------------------------------------------- Set (X,Y) -> 0
func (m *Grid) Unset(X,Y uint) {
	if X < GRID_W && Y < GRID_H {
		bit := Y * GRID_W + X
		n := bit/8
		off := uint(bit%8)
		m.Bitset[n] &= ^(byte(128)>>off)
	}
}
