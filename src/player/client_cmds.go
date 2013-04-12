package player

import "strings"
import "strconv"
import . "db"
import . "types"

func exec_cli(ud *User, msg string) string {
	params:= strings.SplitN(msg, " ", 2)

	switch params[0] {
	case "echo":  return c_echo(ud, params[1])
	case "login": return c_login(ud, params[1])
	case "attack": return c_attack(ud, params[1])
	case "talk": return c_talk(ud, params[1])
	}

	return "Invalid Command"
}

// commands from client
func c_login(ud *User, p string) string {
	ch := make(chan string)
	params:= strings.SplitN(p, " ", 2)

	if len(params) == 2 {
		go DB.Login(ch, params[0], params[1], ud)
		ret := <-ch

		if (ret == "true") {
			RegisterChannel(ud.MQ, ud.Id)
		}
		return ret
	}

	return "false"
}

func c_echo(ud *User, p string) string {
	return p
}

func c_talk(ud *User, p string) string {
	params:= strings.SplitN(p, " ", 2)

	if len(params) >= 2 {
		id,_ := strconv.Atoi(params[0])
		ch := QueryChannel(id)
		if ch != nil {
			msg := []string{"MESG", string(ud.Id), params[1]}
			ch <- strings.Join(msg, " ")
		}
	}

	return "MSG SENT"
}

func c_attack(ud *User, p string) string {
	params:= strings.SplitN(p, " ", 2)

	if len(params) >= 2 {
		id,_ := strconv.Atoi(params[0])
		ch := QueryChannel(id)
		if ch != nil {
			msg := []string{"ATTACKED", string(ud.Id), params[1]}
			ch <- strings.Join(msg, " ")
		}
	}

	return "ATTACK SENT"
}
