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
	var tblname string

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
			tblname = names[0]
			isLineOne = false
			continue
		}

		// the first column is indexed
		for i := 1; i < len(fields); i++ {
			_set(tblname, fields[0], names[i], fields[i])
		}
	}
}
