package types

import (
	"time"
)

/*
 * User is Stored in both ranklist & session.User for efficiency
 * MAKE SURE UPDATE BOTH RankList & Session
 */
type User struct {
	Id           int32
	Name         string
	Mac          string
	Score        int32
	Archives     string
	LastSaveTime time.Time
	ProtectTime  time.Time
	CreatedAt    time.Time
}
