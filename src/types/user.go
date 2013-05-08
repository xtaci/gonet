package types

import (
	"time"
)

type User struct {
	Id           int32
	Name         string
	Mac          string
	Score        int32
	Archives     string
	Bitmap       string
	Buildings    string
	LastSaveTime time.Time
	CreatedAt    time.Time
}
