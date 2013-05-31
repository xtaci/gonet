package gamedata

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	fmt.Println("读取 部落大厅等级 LEVEL:20 Field:金库")
	fmt.Println(GetInt("部落大厅等级", "20", "金库"))
	fmt.Println(GetString("建筑规格", "锻造室", "占用格子数"))
}
