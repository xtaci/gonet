package AI

import (
	"strconv"
	"strings"
)

import (
	"gamedata"
	"types/estates"
	"types/grid"
)

//------------------------------------------------ 根据建筑建立格子
func CreateGrid(manager *estates.Manager) *grid.Grid {
	// 建立位图的格子信息
	Grid := grid.New()
	for k, v := range manager.Estates {
		// TODO :  读gamedata,建立grid信息
		name := gamedata.Query(v.TYPE)
		cell := gamedata.GetString("建筑规格", name, "占用格子数")
		wh := strings.Split(cell, "X")
		w, _ := strconv.Atoi(wh[0])
		h, _ := strconv.Atoi(wh[1])

		oid, _ := strconv.Atoi(k)
		//	fmt.Println(w,h, wh, cell, v)
		for x := int(v.X); x < int(v.X)+w; x++ {
			for y := int(v.Y); y < int(v.Y)+h; y++ {
				Grid.Set(x, y, uint16(oid))
			}
		}
	}

	return Grid
}
