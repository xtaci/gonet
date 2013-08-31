package inspect

import (
	"bufio"
	"fmt"
	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
	"log"
	"net"
	"os"
)

import (
	"agent/gsdb"
)

//---------------------------------------------------------- bind 8800 to localhost
func StartInspect() {
	// Listen
	service := "127.0.0.1:8800"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Fprintln(conn, x)
		}
		conn.Close()
	}()

	fmt.Fprintln(conn, "GameServer LUA Console")

	vm := _lua_vm(conn)
	defer func() {
		vm.Close()
	}()

	// using scanner to read
	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanLines)

	prompt(conn)
	for scanner.Scan() {
		err := vm.DoString(scanner.Text())
		if err != nil {
			fmt.Fprintln(conn, err)
		}
		prompt(conn)
	}
}

func checkError(err error) {
	if err != nil {
		log.Printf("Fatal error: %v", err)
		os.Exit(-1)
	}
}

var script = `function print(...)
				for i=1,arg.n do
					Print(conn, tostring(arg[i]))
				end			
			  end
		      `

func _lua_vm(conn net.Conn) *lua.State {
	lua := luar.Init()
	lua.DoString(script)
	_push_funcs(conn, lua)
	return lua
}

func prompt(conn net.Conn) {
	fmt.Fprint(conn, "> ")
}

//---------------------------------------------------------- funcs need be pushed in
func _push_funcs(conn net.Conn, L *lua.State) {
	luar.Register(L, "", luar.Map{
		"conn":  conn,
		"Print": Print,
	})

	luar.Register(L, "gsdb", luar.Map{
		"ListAll":     gsdb.ListAll,
		"QueryOnline": gsdb.QueryOnline,
	})
}

//---------------------------------------------------------- hijack lua print
func Print(conn net.Conn, s string) {
	fmt.Fprintln(conn, s)
}
