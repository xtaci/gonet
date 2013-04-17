package types

import "time"

type User struct {
	MQ					chan string
	Id                 int
	Name               string
	CreatedAt			time.Time

	Cities []City
}
