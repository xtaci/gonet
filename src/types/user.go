package types

type User struct {
	MQ					chan string
	Id                 int
	Name               string

	Cities []City
}
