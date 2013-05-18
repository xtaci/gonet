package main

import (
	"db/user_tbl"
	"hub/accounts"
)

func startup_work() {
	load_ranklist()
}

//----------------------------------------------- load user table into memory
func load_ranklist() {
	uds := user_tbl.LoadAll()

	for i := range uds {
		accounts.AddUser(&uds[i])
	}

	go accounts.ExpireRoutine()
}
