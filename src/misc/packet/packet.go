package packet

import (
	"errors"
	"math"
)

const (
	PACKET_LIMIT = 65533 // 2^16 - 1 - 2
)

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

func (p *Packet) Seek(n uint) {
	p.pos += n
}

//=============================================== Readers
func (p *Packet) ReadBool() (ret bool, err error) {
	b, _err := p.ReadByte()

	if b != byte(1) {
		return false, _err
	}

	return true, _err
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

func (p *Packet) ReadBytes() (ret []byte, err error) {
	if p.pos+2 > uint(len(p.data)) {
		err = errors.New("read bytes header failed")
		return
	}
	size, _ := p.ReadU16()
	if p.pos+uint(size) > uint(len(p.data)) {
		err = errors.New("read bytes data failed")
		return
	}

	ret = p.data[p.pos : p.pos+uint(size)]
	p.pos += uint(size)
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
	ret = uint16(buf[0])<<8 | uint16(buf[1])
	p.pos += 2
	return
}

func (p *Packet) ReadS16() (ret int16, err error) {
	_ret, _err := p.ReadU16()
	ret = int16(_ret)
	err = _err
	return
}

func (p *Packet) ReadU24() (ret uint32, err error) {
	if p.pos+3 > uint(len(p.data)) {
		err = errors.New("read uint24 failed")
		return
	}

	buf := p.data[p.pos : p.pos+3]
	ret = uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2])
	p.pos += 3
	return
}

func (p *Packet) ReadS24() (ret int32, err error) {
	_ret, _err := p.ReadU24()
	ret = int32(_ret)
	err = _err
	return
}

func (p *Packet) ReadU32() (ret uint32, err error) {
	if p.pos+4 > uint(len(p.data)) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := p.data[p.pos : p.pos+4]
	ret = uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	p.pos += 4
	return
}

func (p *Packet) ReadS32() (ret int32, err error) {
	_ret, _err := p.ReadU32()
	ret = int32(_ret)
	err = _err
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
		ret |= uint64(v) << uint((7-i)*8)
	}
	p.pos += 8
	return
}

func (p *Packet) ReadS64() (ret int64, err error) {
	_ret, _err := p.ReadU64()
	ret = int64(_ret)
	err = _err
	return
}

func (p *Packet) ReadFloat32() (ret float32, err error) {
	bits, _err := p.ReadU32()
	if _err != nil {
		return float32(0), _err
	}

	ret = math.Float32frombits(bits)
	if math.IsNaN(float64(ret)) || math.IsInf(float64(ret), 0) {
		return 0, nil
	}

	return ret, nil
}

func (p *Packet) ReadFloat64() (ret float64, err error) {
	bits, _err := p.ReadU64()
	if _err != nil {
		return float64(0), _err
	}

	ret = math.Float64frombits(bits)
	if math.IsNaN(ret) || math.IsInf(ret, 0) {
		return 0, nil
	}

	return ret, nil
}

//================================================ Writers
func (p *Packet) WriteZeros(n int) {
	zeros := make([]byte, n)
	p.data = append(p.data, zeros...)
}

func (p *Packet) WriteBool(v bool) {
	if v {
		p.data = append(p.data, byte(1))
	} else {
		p.data = append(p.data, byte(0))
	}
}

func (p *Packet) WriteByte(v byte) {
	p.data = append(p.data, v)
}

func (p *Packet) WriteBytes(v []byte) {
	p.WriteU16(uint16(len(v)))
	p.data = append(p.data, v...)
}

func (p *Packet) WriteRawBytes(v []byte) {
	p.data = append(p.data, v...)
}

func (p *Packet) WriteString(v string) {
	bytes := []byte(v)
	p.WriteU16(uint16(len(bytes)))
	p.data = append(p.data, bytes...)
}

func (p *Packet) WriteU16(v uint16) {
	buf := make([]byte, 2)
	buf[0] = byte(v >> 8)
	buf[1] = byte(v)
	p.data = append(p.data, buf...)
}

func (p *Packet) WriteS16(v int16) {
	p.WriteU16(uint16(v))
}

func (p *Packet) WriteU24(v uint32) {
	buf := make([]byte, 3)
	buf[0] = byte(v >> 16)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v)
	p.data = append(p.data, buf...)
}

func (p *Packet) WriteU32(v uint32) {
	buf := make([]byte, 4)
	buf[0] = byte(v >> 24)
	buf[1] = byte(v >> 16)
	buf[2] = byte(v >> 8)
	buf[3] = byte(v)
	p.data = append(p.data, buf...)
}

func (p *Packet) WriteS32(v int32) {
	p.WriteU32(uint32(v))
}

func (p *Packet) WriteU64(v uint64) {
	buf := make([]byte, 8)
	for i := range buf {
		buf[i] = byte(v >> uint((7-i)*8))
	}

	p.data = append(p.data, buf...)
}

func (p *Packet) WriteS64(v int64) {
	p.WriteU64(uint64(v))
}

func (p *Packet) WriteFloat32(f float32) {
	v := math.Float32bits(f)
	p.WriteU32(v)
}

func (p *Packet) WriteFloat64(f float64) {
	v := math.Float64bits(f)
	p.WriteU64(v)
}

func Reader(data []byte) *Packet {
	return &Packet{pos: 0, data: data}
}

func Writer() *Packet {
	pkt := &Packet{pos: 0}
	pkt.data = make([]byte, 0, 128)
	return pkt
}
