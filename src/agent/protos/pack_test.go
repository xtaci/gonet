package protos

import "testing"
import "fmt"
import "misc/packet"

type SUB struct {
	H int16
	I uint16
}
type TEST struct {
	A int32
	B string
	C float32
	D uint32
	F []byte
	Sub []SUB
	m int
}

func TestPack(t *testing.T) {

	test := TEST{A:1, B:"1", C:1.0, D:10, F:[]byte{1,2,3,4,5}}
	test.Sub = make([]SUB,2)
	test.Sub[0].H = 1024
	test.Sub[0].I = 2048
	test.Sub[1].H = 4096
	test.Sub[1].I = 8192

	writer := packet.Writer()
	pack(test, writer)

	fmt.Println(writer.Data())
}
