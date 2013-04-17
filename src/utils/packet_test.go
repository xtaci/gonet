package utils

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

	if a != reader.ReadByte() {
			t.Error("packet readbyte mismatch")
	}

	if b != reader.ReadU16() {
			t.Error("packet readu16 mismatch")
	}

	if c != reader.ReadU24() {
			t.Error("packet readu24 mismatch")
	}

	if d != reader.ReadU32() {
			t.Error("packet readu32 mismatch")
	}

	if e != reader.ReadU64() {
		t.Error("packet readu64 mismatch")
	}

	str := reader.ReadString()
	t.Log(str)
	if "hello world" != str {
		t.Error("packet read string mistmatch")
	}
}
