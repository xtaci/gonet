package gamedata

import (
	"bufio"
	"os"
	"strings"
)

//----------------------------------------------- parse & load a game data file into dictionary
func parse(file *os.File) {
	r := bufio.NewReader(file)

	isLineOne := true
	var names []string

	for {
		line, e := r.ReadString('\n')
		line = strings.TrimSpace(line)

		// empty-line & #comment
		if line == "" {
			if e == nil {
				continue
			} else {
				break
			}
		}

		// split fields
		fields := strings.Split(line, ",")

		// the first line represents field names
		if isLineOne {
			names = make([]string, len(fields))
			for k, v := range fields {
				names[k] = v
			}
			isLineOne = false
			continue
		}

		// 第一列包含特殊含义
		for i := 1; i < len(fields); i++ {
			Set(names[0], fields[0], names[i], fields[i])
		}
	}
}
