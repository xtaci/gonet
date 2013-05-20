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
	Encode(ctx, data)
	fmt.Println(data)
}
