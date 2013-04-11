package player

import "strings"
import "db"


// commands from client
func Client_login(p string) string {
	ch := make(chan string)
	params:= strings.SplitN(p, " ", 2)

	if len(params) == 2 {
		go db.Login(ch, params[0], params[1])
		ret := <-ch

		if (ret == "true") {
			user.name = params[0]
			RegisterChannel(mq, params[0])
		}
		return ret
	}

	return "false"
}

func Client_echo(p string) string {
	return p
}

func Client_talk(p string) string {
	params:= strings.SplitN(p, " ", 2)

	if len(params) >= 2 {
		ch := QueryChannel(params[0])
		if ch != nil {
			print(params[1])
			msg := []string{"mesg", params[1]}
			ch <- strings.Join(msg, " ")
		}
	}

	return ""
}

func Client_attack(p string) string {
	ch := QueryChannel(p)

	if ch != nil {
		msg := []string{"attackedby", user.name}
		ch <- strings.Join(msg, " ")
	}

	return ""
}
