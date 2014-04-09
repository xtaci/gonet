package stats

import (
	"time"
)

const (
	INT_GAME_INFO = "INT_GAME_INFO"
	STR_GAME_INFO = "STR_GAME_INFO"
	PAYS_INFO     = "PAYS_INFO"
	LTV_INFO      = "LTV_INFO"
	NEW_USER      = "NEW_USER"
)

type IntGameInfo struct {
	IntValue int32
	Key      string
	Time     time.Time
	Lang     string
}

type StrGameInfo struct {
	StrValue string
	Key      string
	Time     time.Time
	Lang     string
}

//----------------------------------------------- 玩家付费信息
// PAYS_INFO
type Pays struct {
	Time        time.Time
	Device_id   string
	Appid       string
	Server_name string
	Unix_time   int64
	Money       float64
	Gems        int32
	Id          int32
}

//----------------------------------------------- 留存
// LTV_INFO
type LTV struct {
	Time      time.Time
	Pays      float64 // 新用户付费数
	Res       int64   // 新用户留存
	Sign_time int64   // 注册时间
	Cnt       int64   // 新用户注册总数
	Lang      string  // 语言
}

//----------------------------------------------- 新用户总数
// NEW_USER
type New struct {
	Sign_time int64 // 注册时间 (统一为每天的0点)
	Cnt       int64 // 新用户注册总数
}
