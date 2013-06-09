package protos

import (
	"fmt"
	"log"
	"misc/packet"
	"net"
	"os"
	"testing"
)

func TestAgent(t *testing.T) {
	log.Println("Connecting to GS")
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	U := user_login_info{}
	pkt := packet.Pack(Code["user_login_req"], U, nil)
	fmt.Println(pkt)

	writer := packet.Writer()
	writer.WriteU16(uint16(len(pkt) + 4))
	writer.WriteU32(0)
	writer.WriteRawBytes(pkt)

	fmt.Println(writer.Data())
	conn.Write(writer.Data())

	ret := make([]byte, 100)
	conn.Read(ret)
	fmt.Println(ret)
}

func BenchmarkAgent(b *testing.B) {
	log.Println("Connecting to GS")
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	U := user_login_info{}
	pkt := packet.Pack(Code["user_login_req"], U, nil)

	writer := packet.Writer()
	writer.WriteU16(uint16(len(pkt) + 4))
	writer.WriteU32(0)
	writer.WriteRawBytes(pkt)
	ret := make([]byte, 100)

	fmt.Println("Benchmark", b.N)
	for i := 0; i < b.N; i++ {
		conn.Write(writer.Data())
		conn.Read(ret)
	}
}
