package packet

import "testing"

func TestPacketWriter(t *testing.T) {
	p := PacketWriter()
	a := byte(0xFF)
	b := uint16(0xFF00)
	c := uint32(0xFF0000)
	d := uint32(0xFF000000)
	e := uint64(0xFF00000000000000)

	p.WriteByte(a)
	p.WriteU16(b)
	p.WriteU24(c)
	p.WriteU32(d)
	p.WriteU64(e)

	p.WriteString("hello world")

	reader := PacketReader(p.Data())

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

	tmp5, _ := reader.ReadString()

	if "hello world" != tmp5 {
		t.Error("packet read string mistmatch")
	}

	_, err := reader.ReadByte()

	if err == nil {
		t.Error("overflow check failed")
	}
}
