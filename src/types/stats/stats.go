package stats

import (
	"time"
)

const (
	INT_GAME_INFO = "INT_GAME_INFO"
	STR_GAME_INFO = "STR_GAME_INFO"
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
