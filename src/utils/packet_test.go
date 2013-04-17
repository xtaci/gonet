package utils

import "testing"

func TestPacketWriter(t *testing.T) {
	p := PacketWriter()
	a := byte(0xFF)
	b := uint16(0xFF00)
	c := uint32(0xFF0000)
	d := uint32(0xFF000000)

	p.WriteByte(a)
	p.WriteU16(b)
	p.WriteU24(c)
	p.WriteU32(d)

	data := p.Data()
	result := []byte{255,0,255,0,0,255,0,0,0,255}

	for i := range data {
		if result[i] != data[i] {
			t.Error("packet writer failed")
		}
	}
}
