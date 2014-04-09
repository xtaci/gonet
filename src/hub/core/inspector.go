package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
	"log"
	"net"
	"os"
)

func StartInspect() {
	var conn net.Conn
	// Listen
	service := "127.0.0.1:8801"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err = listener.Accept()

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

	fmt.Fprintln(conn, "HUB LUA Console")
	fmt.Fprintln(conn, "eg:")
	fmt.Fprintln(conn, "\tprintusers()")
	fmt.Fprintln(conn, "\tprintjson(players[1])")

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

              function printjson(x)
                print(goprintjson(x))
              end
				
			  function printusers()
				table.foreach(luar.map2table(players), print)
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
		"players":     _players,
		"conn":        conn,
		"Print":       Print,
		"goprintjson": Printjson,
	})
}

func Print(conn net.Conn, s string) {
	fmt.Fprintln(conn, s)
}

//---------------------------------------------------------- printjson
func Printjson(x interface{}) string {
	b, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(b)
}
