package gamedata

import (
	"bufio"
	"os"
	"strings"
)

//----------------------------------------------- parse & load a game data file into dictionary
func parse(file *os.File) {
	isLineOne := true
	var names []string
	var tblname string

	// using scanner to read csv file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// ignore empty-line
		if line == "" {
			continue
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
