package packet

import (
	"fmt"
	"testing"
)

func TestPacketWriter(t *testing.T) {
	p := Writer()
	a := byte(0xFF)
	b := uint16(0xFF00)
	c := uint32(0xFF0000)
	d := uint32(0xFF000000)
	e := uint64(0xFF00000000000000)
	f32 := float32(1.0)
	f64 := float64(2.0)

	p.WriteBool(true)
	p.WriteByte(a)
	p.WriteU16(b)
	p.WriteU24(c)
	p.WriteU32(d)
	p.WriteU64(e)
	p.WriteFloat32(f32)
	p.WriteFloat64(f64)

	p.WriteString("hello world")
	p.WriteBytes([]byte("hello world"))
	var nil_bytes []byte
	p.WriteBytes(nil_bytes)

	reader := Reader(p.Data())

	BOOL, _ := reader.ReadBool()
	if BOOL != true {
		t.Error("packet readbool mismatch")
	}

	tmp, _ := reader.ReadByte()
	if a != tmp {
		t.Error("packet readbyte mismatch")
	}

	tmp1, _ := reader.ReadU16()
	if b != tmp1 {
		t.Error("packet readu16 mismatch")
	}

	tmp2, _ := reader.ReadU24()
	if c != tmp2 {
		t.Error("packet readu24 mismatch")
	}

	tmp3, _ := reader.ReadU32()
	if d != tmp3 {
		t.Error("packet readu32 mismatch")
	}

	tmp4, _ := reader.ReadU64()
	if e != tmp4 {
		t.Error("packet readu64 mismatch")
	}

	tmp5, _ := reader.ReadFloat32()
	if f32 != tmp5 {
		t.Error("packet readf32 mismatch")
	}

	tmp6, _ := reader.ReadFloat64()
	if f64 != tmp6 {
		t.Error("packet readf32 mismatch")
	}

	tmp100, _ := reader.ReadString()

	if "hello world" != tmp100 {
		t.Error("packet read string mistmatch")
	}

	tmp101, _ := reader.ReadBytes()

	fmt.Println(tmp101)
	if tmp101[0] != 'h' {
		t.Error("packet read bytes mistmatch")
	}

	tmp102, _ := reader.ReadBytes()
	fmt.Println("NIL:", tmp102)
	if len(tmp102) != 0 {
		t.Error("packet read nil bytes mistmatch")
	}

	_, err := reader.ReadByte()

	if err == nil {
		t.Error("overflow check failed")
	}
}

func BenchmarkPacketWriter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := Writer()
		p.WriteU16(128)
		p.WriteBool(true)
		p.WriteS32(-16)
		p.WriteString("A")
		p.WriteFloat32(1.0)
		p.WriteU32(16)
		p.WriteFloat64(1.0)
	}
}
