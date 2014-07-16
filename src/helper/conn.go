package helper

import (
	"net"
)

func SetConnParam(conn *net.TCPConn) {
	conn.SetNoDelay(false)
	conn.SetKeepAlive(true)
	conn.SetLinger(-1)
}
