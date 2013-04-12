package cmd

import "strings"
import "strconv"
import . "db"
import . "types"
import "names"

type ClientCmds struct {
}

func ExecCli(ud *User, msg string) string {
	var cmd ClientCmds;

	params := strings.SplitN(msg, " ", 2)

	switch params[0] {
	case "echo":
		return cmd.echo(ud, params[1])
	case "login":
		return cmd.login(ud, params[1])
	case "attack":
		return cmd.attack(ud, params[1])
	case "talk":
		return cmd.talk(ud, params[1])
	}

	return "Invalid Command"
}

// commands from client
func (ClientCmds) login(ud *User, p string) string {
	ch := make(chan string)
	params := strings.SplitN(p, " ", 2)

	if len(params) == 2 {
		go DB.Login(ch, params[0], params[1], ud)
		ret := <-ch

		if ret == "true" {
			names.RegisterChannel(ud.MQ, ud.Id)
		}
		return ret
	}

	return "false"
}

func (ClientCmds) echo(ud *User, p string) string {
	return p
}

func (ClientCmds) talk(ud *User, p string) string {
	params := strings.SplitN(p, " ", 2)

	if len(params) >= 2 {
		id, _ := strconv.Atoi(params[0])
		ch := names.QueryChannel(id)
		if ch != nil {
			msg := []string{"MESG", string(ud.Id), params[1]}
			ch <- strings.Join(msg, " ")
		}
	}

	return "MSG SENT"
}

func (ClientCmds) attack(ud *User, p string) string {
	params := strings.SplitN(p, " ", 2)

	if len(params) >= 2 {
		id, _ := strconv.Atoi(params[0])
		ch := names.QueryChannel(id)
		if ch != nil {
			msg := []string{"ATTACKED", string(ud.Id), params[1]}
			ch <- strings.Join(msg, " ")
		}
	}

	return "ATTACK SENT"
}
