package packet

import "testing"
import "fmt"

type SUB2 struct {
	M []int16
}

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
	S2   SUB
	S3   SUB2
}

type TEST2 struct {
	BOOL bool
	A    int32
	B    string
	C    float32
	D    uint32
	E    float64
}

func TestPack(t *testing.T) {

	test := TEST{BOOL: true, A: 16, B: string([]byte{65}), C: 1.0, D: 32, E: 1.0, F: []byte{1, 2, 3, 4, 5}}
	test.Sub = make([]SUB, 2)
	test.Sub[0].H = 1024
	test.Sub[0].I = 2048
	test.Sub[1].H = 4096
	test.Sub[1].I = 8192
	test.S2.H = 100
	test.S3.M = make([]int16, 10)

	fmt.Println(Pack(128, test, nil))
	fmt.Println(Pack(128, &test, nil))
	fmt.Println(Pack(129, nil, nil))
}

func BenchmarkPack(b *testing.B) {
	test := TEST2{BOOL: true, A: 16, B: string([]byte{65}), C: 1.0, D: 32, E: 1.0}
	for i := 0; i < b.N; i++ {
		Pack(128, test, nil)
	}
}
