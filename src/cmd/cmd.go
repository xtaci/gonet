package cmd

import "strings"
import . "types"
import cli "cmd/cli"
import srv "cmd/srv"
//import "log"

func ExecCli(ud *User, p []byte) string {
	switch p[0] {
	case 'E':
		return cli.Echo(ud, string(p[1:]))
	case 'L':
		return cli.Login(ud, string(p[1:]))
	case 'A':
		return cli.Attack(ud, string(p[1:]))
	case 'T':
		return cli.Talk(ud, string(p[1:]))
	case 'N':
		return cli.Newcity(ud, string(p[1:]))
	}

	return "Invalid Command"
}

func ExecSrv(ud *User, msg string) string {
	params := strings.SplitN(msg, " ", 2)
	switch params[0] {
	case "MESG":
		return srv.Mesg(ud, params[1])
	case "ATTACKED":
		return srv.Attacked(ud, params[1])
	}

	return ""
}
