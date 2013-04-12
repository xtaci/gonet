package main

import "os"
import "io"
import "bufio"
import "regexp"
import "fmt"

func read_config(path string)(ret map[string]string) {
	ret = make(map[string]string)
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("error opening file %v\n", err)
		os.Exit(1)
	}

	re := regexp.MustCompile(`[\t ]*([0-9A-Za-z_]+)[\t ]*=[\t ]*([0-9A-Za-z_:?]+)[\t ]*`)

	r := bufio.NewReader(f)
	for {
		line,e := r.ReadString('\n')

		if line == "" || []byte(line)[0] == '#' {
			continue
		}
		// maping
		slice := re.FindStringSubmatch(line)

		if slice != nil{
			ret[slice[1]] = slice[2]
		}

		if e == io.EOF {
			break
		}
	}

	return
}
