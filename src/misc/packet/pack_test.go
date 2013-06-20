package packet

import "testing"
import "fmt"

type SUB struct {
	H int16
	I uint16
}
type TEST struct {
	BOOL bool
	A    int32
	B    string
	C    float32
	D    uint32
	E    float64
	F    []byte
	Sub  []SUB
}

func TestPack(t *testing.T) {

	test := TEST{BOOL: true, A: 16, B: string([]byte{65}), C: 1.0, D: 32, E: 1.0, F: []byte{1, 2, 3, 4, 5}}
	test.Sub = make([]SUB, 2)
	test.Sub[0].H = 1024
	test.Sub[0].I = 2048
	test.Sub[1].H = 4096
	test.Sub[1].I = 8192

	fmt.Println(Pack(128, test, nil))
	fmt.Println(Pack(128, &test, nil))
}
