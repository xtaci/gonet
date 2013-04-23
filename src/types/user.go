package types

import "time"

type User struct {
	Id        int
	Name      string
	CreatedAt time.Time
}
