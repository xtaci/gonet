package utils

type Packet struct {
	pos uint16
	data []byte
}

func (p *Packet) Data() []byte {
	return p.data
}

//---------------------------------------------------------Reader
func (p *Packet) SkipN(n int) {
	p.pos+=uint16(n)
}

func (p *Packet) ReadByte() (ret byte) {
	ret = p.data[p.pos]
	p.pos++
	return
}

func (p *Packet) ReadU16() (ret uint16){
	buf := p.data[p.pos:p.pos+2]
	ret = uint16(buf[1])<<8|uint16(buf[0])
	p.pos+=2
	return
}

func (p *Packet) ReadU24() (ret uint32) {
	buf := p.data[p.pos:p.pos+3]
	ret = uint32(buf[2])<<8 | uint32(buf[1])<<8 | uint32(buf[0])
	p.pos+=3
	return
}

func (p *Packet) ReadU32() (ret uint32) {
	buf := p.data[p.pos:p.pos+4]
	ret = uint32(buf[3])<<8 | uint32(buf[2])<<8 | uint32(buf[1])<<8 | uint32(buf[0])
	p.pos+=4
	return
}

func (p *Packet) ReadU64() (ret uint64) {
	ret=0
	buf := p.data[p.pos:p.pos+8]
	for i, v := range buf {
		ret |= uint64(v) << uint(i*8)
	}
	p.pos+=8
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

func (p *Packet) WriteU16(v uint16) {
	p.data = append(p.data, byte(v))
	p.data = append(p.data, byte(v>>8))
}

func (p *Packet) WriteU24(v uint32) {
	p.data = append(p.data, byte(v))
	p.data = append(p.data, byte(v>>8))
	p.data = append(p.data, byte(v>>16))
}

func (p *Packet) WriteU32(v uint32) {
	p.data = append(p.data, byte(v))
	p.data = append(p.data, byte(v>>8))
	p.data = append(p.data, byte(v>>16))
	p.data = append(p.data, byte(v>>24))
}


func PacketReader(data []byte) *Packet{
	return &Packet{pos:0, data:data}
}

func PacketWriter() *Packet {
	return &Packet{pos:0}
}
