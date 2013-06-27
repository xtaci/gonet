package ipc

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

	U := &user_login_info{}
	U.F_user_name = "test1"
	U.F_mac_addr = "mac1"
	pkt := packet.Pack(Code["user_login_req"], U, nil)

	writer := packet.Writer()
	writer.WriteU16(uint16(len(pkt) + 4))
	writer.WriteU32(0)
	writer.WriteRawBytes(pkt)

	conn.Write(writer.Data())

	ret := make([]byte, 100)
	n, _ := conn.Read(ret)
	fmt.Printf("%q\n", ret[:n])

	// talk
	msg := &talk{}
	msg.F_user = "test1"
	msg.F_msg = "hello world"
	pkt = packet.Pack(Code["talk_req"], msg, nil)

	writer = packet.Writer()
	writer.WriteU16(uint16(len(pkt) + 4))
	writer.WriteU32(0)
	writer.WriteRawBytes(pkt)

	conn.Write(writer.Data())
	n, _ = conn.Read(ret)
	fmt.Printf("%q\n", ret[:n])
}

func BenchmarkAgent(b *testing.B) {
	fmt.Println("Benchmark", b.N)
	for i := 0; i < b.N; i++ {
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
		U.F_user_name = fmt.Sprintf("test%v", i)
		U.F_mac_addr = fmt.Sprintf("mac%v", i)

		pkt := packet.Pack(Code["user_login_req"], &U, nil)

		writer := packet.Writer()
		writer.WriteU16(uint16(len(pkt) + 4))
		writer.WriteU32(0)
		writer.WriteRawBytes(pkt)
		ret := make([]byte, 100)

		conn.Write(writer.Data())
		conn.Read(ret)
	}
}
