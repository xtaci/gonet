package cmd

import "strings"
import "math/rand"
import . "types"

type ServerCmds struct {
}

func ExecSrv(ud *User, msg string) string {
	var cmds ServerCmds
	params := strings.SplitN(msg, " ", 2)
	switch params[0] {
	case "MESG":
		return cmds.mesg(ud, params[1])
	case "ATTACKED":
		return cmds.attacked(ud, params[1])
	}

	return ""
}

func (ServerCmds) mesg(ud *User, p string) string {
	msg := []string{"mesg", p}
	return strings.Join(msg, " ")
}

func (ServerCmds) attacked(ud *User, p string) string {
	msg := []string{"attacked", p}

	//
	nums := make([]int, 100)
	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Int()
	}

	// insertion sort
	for i := 1; i < len(nums); i++ {
		cur := nums[i]
		j := i - 1

		for j >= 0 && nums[j] > cur {
			nums[j+1] = nums[j]
			j = j - 1
		}

		nums[j+1] = cur
	}

	/*
		for i:=1; i<len(nums);i++ {
			println(nums[i])
		}
	*/

	return strings.Join(msg, " ")
}
