package main

import (
	"db/user_tbl"
	"hub/core"
)

func startup_work() {
	load_ranklist()
}

//----------------------------------------------- load user table into memory
func load_ranklist() {
	// load users
	uds := user_tbl.GetAll()

	for i := range uds {
		core.AddUser(&uds[i])
	}
}
