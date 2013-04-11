package main

import "net"
import "player"
import "io"

func main() {
	println("Starting the server")

	listener, err := net.Listen("tcp", "0.0.0.0:8888")

	if err != nil {
		println("error listening:", err.Error())
		return
	}


	for {
		conn, err := listener.Accept()
		if err != nil {
			println("error accept::", err.Error())
			return
		}
		go HandleClient(conn)
	}
}

func HandleClient(conn net.Conn) {
	header := make([]byte, 2)
	ch := make(chan string)

	go player.Start(ch, conn)

	for {
		// header
		n, err := io.ReadFull(conn, header)
		if n==0 && err == io.EOF {
			break;
		} else if err != nil {
			println("error receving header:", err)
			break;
		}

		// data
		size := int(header[0] <<8 | header[1])
		println(size)
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err  != nil {
			println("error receving msg:", err)
			break
		}

		ch <- string(data)
	}

	ch <- "CLIENTCLOSE"
	conn.Close()
	close(ch)
}
