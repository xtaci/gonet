package start

import (
	"db/user_tbl"
	"hub/ranklist"
)

func startup_work() {
	load_ranklist()
}

//----------------------------------------------- load user table into memory
func load_ranklist() {
	uds := user_tbl.LoadAll()

	for i := range uds {
		ranklist.AddUser(&uds[i])
	}

	go ranklist.ExpireRoutine()
}
