package player

import "strings"

func (user *UserData) exec_srv(msg string) string {
	params:= strings.SplitN(msg, " ", 2)
	switch params[0] {
	case "MESG": return user.S_mesg(params[1]);
	case "ATTACKED": return user.S_attacked(params[1]);
	}

	return ""
}

func (user *UserData) S_mesg(p string) string {
	msg := []string{"mesg",p}
	return strings.Join(msg, " ")
}

func (user *UserData) S_attacked(p string) string {
	return "attacked by" +  p
}
