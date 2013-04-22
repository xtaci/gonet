package cmd

import "strings"
import "math/rand"
import . "types"

func Mesg(ud *User, p string) string {
	msg := []string{"mesg", p}
	return strings.Join(msg, " ")
}

func Attacked(ud *User, p string) string {
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
