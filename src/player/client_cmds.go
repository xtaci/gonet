package player

import "strings"
import "db"

func (user *UserData) exec_cli(msg string) string {
	params:= strings.SplitN(msg, " ", 2)

	switch params[0] {
	case "echo":  return user.C_echo(params[1])
	case "login": return user.C_login(params[1])
	case "attack": return user.C_attack(params[1])
	case "talk": return user.C_talk(params[1])
	}

	return "Invalid Command"
}

// commands from client
func (user *UserData) C_login(p string) string {
	ch := make(chan string)
	params:= strings.SplitN(p, " ", 2)

	if len(params) == 2 {
		go db.Login(ch, params[0], params[1])
		ret := <-ch

		if (ret == "true") {
			user.name = params[0]
			RegisterChannel(user.mq, params[0])
		}
		return ret
	}

	return "false"
}

func (user *UserData) C_echo(p string) string {
	return p
}

func (user *UserData ) C_talk(p string) string {
	params:= strings.SplitN(p, " ", 2)

	if len(params) >= 2 {
		ch := QueryChannel(params[0])
		if ch != nil {
			msg := []string{"MESG", user.name, params[1]}
			ch <- strings.Join(msg, " ")
		}
	}

	return ""
}

func (user *UserData) C_attack(p string) string {
	ch := QueryChannel(p)

	if ch != nil {
		msg := []string{"ATTACKED", user.name}
		ch <- strings.Join(msg, " ")
	}

	return ""
}
