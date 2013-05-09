package gamedata

import (
	"bufio"
	//	"fmt"
	"os"
	"strconv"
	"strings"
)

func parse(file *os.File) {
	r := bufio.NewReader(file)

	isLineOne := true
	var names []string

	for {
		line, e := r.ReadString('\n')
		line = strings.TrimSpace(line)

		if e != nil {
			break
		}

		// split fields
		fields := strings.Split(line, ",")

		if isLineOne {
			names = make([]string, len(fields))
			for k, v := range fields {
				names[k] = v
			}
			isLineOne = false
			continue
		}

		// the first column is LEVEL
		lv, _ := strconv.Atoi(fields[0])

		for i := 1; i < len(fields); i++ {
			Set(names[0], lv, names[i], fields[i])
			//	fmt.Printf("%v %v %v %v\n", names[0],lv, names[i], fields[i])
		}
	}
}
