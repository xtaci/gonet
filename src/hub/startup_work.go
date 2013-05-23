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
	uds := user_tbl.GetAll()

	for i := range uds {
		if uds[i] != nil { // in case of db corruption
			accounts.AddUser(uds[i])
		}
	}
}
