package player

import "strings"
import "strconv"

func (ud *UserData) exec_cli(msg string) string {
	params:= strings.SplitN(msg, " ", 2)

	switch params[0] {
	case "echo":  return ud.C_echo(params[1])
	case "login": return ud.C_login(params[1])
	case "attack": return ud.C_attack(params[1])
	case "talk": return ud.C_talk(params[1])
	}

	return "Invalid Command"
}

// commands from client
func (ud *UserData) C_login(p string) string {
	ch := make(chan string)
	params:= strings.SplitN(p, " ", 2)

	if len(params) == 2 {
		go DBLogin(ch, params[0], params[1], ud)
		ret := <-ch

		if (ret == "true") {
			RegisterChannel(ud.mq, ud.id)
		}
		return ret
	}

	return "false"
}

func (ud *UserData) C_echo(p string) string {
	return p
}

func (ud *UserData) C_talk(p string) string {
	params:= strings.SplitN(p, " ", 2)

	if len(params) >= 2 {
		id,_ := strconv.Atoi(params[0])
		ch := QueryChannel(id)
		if ch != nil {
			msg := []string{"MESG", string(ud.id), params[1]}
			ch <- strings.Join(msg, " ")
		}
	}

	return "MSG SENT"
}

func (ud *UserData) C_attack(p string) string {
	params:= strings.SplitN(p, " ", 2)

	if len(params) >= 2 {
		id,_ := strconv.Atoi(params[0])
		ch := QueryChannel(id)
		if ch != nil {
			msg := []string{"ATTACKED", string(ud.id), params[1]}
			ch <- strings.Join(msg, " ")
		}
	}

	return "ATTACK SENT"
}
