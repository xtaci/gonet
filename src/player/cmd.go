package player

import "strings"
import . "types"
import cli "player/cli"
import srv "player/srv"
import "utils"
import "log"

func ExecCli(ud *User, p []byte) string {
	reader := utils.PacketReader(p)

	b, err := reader.ReadByte()

	if err!=nil {
		log.Println("read protocol error")
	}

	subp := string(reader.Data()[reader.Pos():])

	switch b {
	case 'E':
		return cli.Echo(ud, subp)
	case 'L':
		return cli.Login(ud, subp)
	case 'A':
		return cli.Attack(ud, subp)
	case 'T':
		return cli.Talk(ud, subp)
	case 'N':
		return cli.Newcity(ud, subp)
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
