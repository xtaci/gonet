package stats_client

var Code = map[string]int16{
	"set_adds_req":   100, // 累计信息
	"set_update_req": 200, // 更新信息
}

var RCode = map[int16]string{
	100: "set_adds_req",   // 累计信息
	200: "set_update_req", // 更新信息
}
