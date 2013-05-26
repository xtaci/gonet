package gamedata

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	fmt.Println(FieldNames("部落大厅等级"))
	fmt.Println(NumLevels("部落大厅等级"))
	fmt.Println("读取 部落大厅等级 20 金库")
	fmt.Println(GetInt("部落大厅等级", 20, "金库"))
	fmt.Println(GetFloat("部落大厅等级", 20, "金库"))
}
