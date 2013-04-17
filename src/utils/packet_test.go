package utils

import "testing"

func TestPacketWriter(t *testing.T) {
	p := PacketWriter()
	a := byte(240)
	b := uint16(61680)
	c := uint32(15790320)
	d := uint32(4042322160)

	p.WriteByte(a)
	p.WriteU16(b)
	p.WriteU24(c)
	p.WriteU32(d)

	data := p.Data()
	result := []byte{0xF0,0xF0,0xF0,0xF0,0xF0,0xF0,0xF0, 0xF0, 0xF0, 0xF0}

	for i := range data {
		if result[i] != data[i] {
			t.Error("packet writer failed")
		}
	}
}
