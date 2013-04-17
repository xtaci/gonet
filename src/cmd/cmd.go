package cmd

import "strings"
import . "types"
import cli "cmd/cli"
import srv "cmd/srv"
//import "log"

func ExecCli(ud *User, msg string) string {
	params := strings.SplitN(msg, " ", 2)

	switch params[0] {
	case "echo":
		return cli.Echo(ud, params[1])
	case "login":
		return cli.Login(ud, params[1])
	case "attack":
		return cli.Attack(ud, params[1])
	case "talk":
		return cli.Talk(ud, params[1])
	case "newcity":
		return cli.Newcity(ud, params[1])
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
