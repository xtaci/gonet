package gamedata

import (
	"fmt"
	"misc/naming"
	"testing"
)

func TestParser(t *testing.T) {
	fmt.Println("读取 部落大厅等级 LEVEL:20 Field:金库")
	fmt.Println(GetInt("部落大厅等级", "20", "金库"))
	fmt.Println(GetString("建筑规格", "锻造室", "占用格子数"))

	hash := naming.FNV1a("锻造室")
	fmt.Println(Query(hash))
}

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetString("建筑规格", "锻造室", "占用格子数")
		GetInt("部落大厅等级", "20", "金库")
	}
}
