package pike

import "testing"
import "fmt"

func TestPike(t *testing.T) {
	ctx := NewCtx(1234)
	fmt.Println("###")
	fmt.Println(ctx.sd)
	fmt.Println(ctx.addikey[0].buffer)
	fmt.Println(ctx.addikey[1].buffer)
	fmt.Println(ctx.addikey[2].buffer)

	data := make([]byte, 1024*1024)
	for i := 0; i < len(data); i++ {
		data[i] = byte(i % 256)
	}
	ctx.Codec(data)
	//fmt.Println("ciphertext:", string(data), len(data))
	ctx1 := NewCtx(1234)
	ctx1.Codec(data)
	for i := 0; i < 1024*1024; i++ {
		if data[i] != byte(i%256) {
			t.Error("解码错误")
		}
	}
}
