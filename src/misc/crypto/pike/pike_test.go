package pike

import "testing"
import "fmt"

func TestPike(t *testing.T) {
	ctx := NewCtx(0)
	fmt.Println("###")
	fmt.Println(ctx.sd)
	fmt.Println(ctx.addikey[0].buffer)
	fmt.Println(ctx.addikey[1].buffer)
	fmt.Println(ctx.addikey[2].buffer)
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	ctx.Codec(data)
	fmt.Println("encode:", data)
	ctx1 := NewCtx(0)
	ctx1.Codec(data)
	fmt.Println("decode:", data)
}
