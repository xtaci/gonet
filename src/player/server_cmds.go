package player

import "strings"
import "math/rand"
import . "types"

func exec_srv(ud *User, msg string) string {
	params:= strings.SplitN(msg, " ", 2)
	switch params[0] {
	case "MESG": return s_mesg(ud, params[1]);
	case "ATTACKED": return s_attacked(ud, params[1]);
	}

	return ""
}

func s_mesg(ud *User, p string) string {
	msg := []string{"mesg",p}
	return strings.Join(msg, " ")
}

func s_attacked(ud *User, p string) string {
	msg := []string{"attacked",p}

	//
	nums := make([]int, 100)
	for i:=0;i<len(nums);i++ {
		nums[i] = rand.Int()
	}


	// insertion sort
	for i:=1; i<len(nums);i++ {
		cur := nums[i]
		j:=i-1

		for j>=1 && nums[j] > cur {
			nums[j+1] = nums[j]
			j=j-1
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
