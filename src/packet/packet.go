package packet

import "errors"

type Packet struct {
	pos  uint
	data []byte
}

func (p *Packet) Data() []byte {
	return p.data
}

func (p *Packet) Length() int {
	return len(p.data)
}

func (p *Packet) Pos() uint {
	return p.pos
}

//---------------------------------------------------------Reader
func (p *Packet) Seek(n uint) {
	p.pos += n
}

func (p *Packet) ReadByte() (ret byte, err error) {
	if p.pos >= uint(len(p.data)) {
		err = errors.New("read byte failed")
		return
	}

	ret = p.data[p.pos]
	p.pos++
	return
}

func (p *Packet) ReadString() (ret string, err error) {
	if p.pos+2 > uint(len(p.data)) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := p.ReadU16()
	if p.pos+uint(size) > uint(len(p.data)) {
		err = errors.New("read string data failed")
		return
	}

	bytes := p.data[p.pos : p.pos+uint(size)]
	p.pos += uint(size)
	ret = string(bytes)
	return
}

func (p *Packet) ReadU16() (ret uint16, err error) {
	if p.pos+2 > uint(len(p.data)) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := p.data[p.pos : p.pos+2]
	ret = uint16(buf[1])<<8 | uint16(buf[0])
	p.pos += 2
	return
}

func (p *Packet) ReadU24() (ret uint32, err error) {
	if p.pos+3 > uint(len(p.data)) {
		err = errors.New("read uint24 failed")
		return
	}

	buf := p.data[p.pos : p.pos+3]
	ret = uint32(buf[2])<<16 | uint32(buf[1])<<8 | uint32(buf[0])
	p.pos += 3
	return
}

func (p *Packet) ReadU32() (ret uint32, err error) {
	if p.pos+4 > uint(len(p.data)) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := p.data[p.pos : p.pos+4]
	ret = uint32(buf[3])<<24 | uint32(buf[2])<<16 | uint32(buf[1])<<8 | uint32(buf[0])
	p.pos += 4
	return
}

func (p *Packet) ReadU64() (ret uint64, err error) {
	if p.pos+8 > uint(len(p.data)) {
		err = errors.New("read uint64 failed")
		return
	}

	ret = 0
	buf := p.data[p.pos : p.pos+8]
	for i, v := range buf {
		ret |= uint64(v) << uint(i*8)
	}
	p.pos += 8
	return
}

//---------------------------------------------------------Writer
func (p *Packet) WriteZeros(n int) {
	zeros := make([]byte, n)
	p.data = append(p.data, zeros...)
}

func (p *Packet) WriteByte(v byte) {
	p.data = append(p.data, v)
}

func (p *Packet) WriteString(v string) {
	bytes := []byte(v)
	p.WriteU16(uint16(len(bytes)))
	p.data = append(p.data, bytes...)
}

func (p *Packet) WriteU16(v uint16) {
	buf := make([]byte, 2)
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	p.data = append(p.data, buf...)
}

func (p *Packet) WriteU24(v uint32) {
	buf := make([]byte, 3)
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 16)
	p.data = append(p.data, buf...)
}

func (p *Packet) WriteU32(v uint32) {
	buf := make([]byte, 4)
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 16)
	buf[3] = byte(v >> 24)
	p.data = append(p.data, buf...)
}

func (p *Packet) WriteU64(v uint64) {
	buf := make([]byte, 8)
	for i := range buf {
		buf[i] = byte(v >> uint(i*8))
	}

	p.data = append(p.data, buf...)
}

func PacketReader(data []byte) *Packet {
	return &Packet{pos: 0, data: data}
}

func PacketWriter() *Packet {
	return &Packet{pos: 0}
}
