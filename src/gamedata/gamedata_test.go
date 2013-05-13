package gamedata

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {
	files, _ := filepath.Glob("./data/*.csv")

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			fmt.Println("error opening file %v\n", err)
			continue
		}

		parse(file)
		file.Close()
	}

	fmt.Println("读取 部落大厅等级 20 金库")
	fmt.Println(Get("部落大厅等级", 20, "金库"))
}
