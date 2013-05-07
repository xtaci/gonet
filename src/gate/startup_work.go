package main

import (
	"db/user_tbl"
	"hub/ranklist"
)

// load user table into memory
func load_ranklist() {
	uds := user_tbl.ReadAll()

	for i := range uds {
		ranklist.AddUser(&uds[i])
	}
}
