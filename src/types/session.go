package types

type Session struct {
	MQ		chan string
	User	User
	Cities	[]City
}
