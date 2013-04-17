package cmd

import "strings"
import "strconv"
import "db/user"
import "db/city"
import . "types"
import "names"
//import "log"

// commands from client
func Login(ud *User, p string) string {
	ch := make(chan string)
	params := strings.SplitN(p, " ", 2)

	if len(params) == 2 {
		go user.Login(ch, params[0], params[1], ud)
		ret := <-ch

		if ret == "true" {
			names.Register(ud.MQ, ud.Id)
			// load cities
			ud.Cities = city.LoadCities(ud.Id)
		}
		return ret
	}

	return "false"
}

func Echo(ud *User, p string) string {
	return p
}

func Talk(ud *User, p string) string {
	params := strings.SplitN(p, " ", 2)

	if len(params) >= 2 {
		id, _ := strconv.Atoi(params[0])
		ch := names.Query(id)
		if ch != nil {
			msg := []string{"MESG", string(ud.Id), params[1]}
			ch <- strings.Join(msg, " ")
		}
	}

	return "MSG SENT"
}

func Attack(ud *User, p string) string {
	params := strings.SplitN(p, " ", 2)

	if len(params) >= 2 {
		id, _ := strconv.Atoi(params[0])
		ch := names.Query(id)
		if ch != nil {
			msg := []string{"ATTACKED", string(ud.Id), params[1]}
			ch <- strings.Join(msg, " ")
		}
	}

	return "ATTACK SENT"
}

func Newcity(ud *User, p string) string {
	newcity := City { Name:p, OwnerId:ud.Id }
	ud.Cities = append(ud.Cities, newcity)
	city.Create(&ud.Cities[len(ud.Cities)-1])
	return "CITY CREATED"
}
