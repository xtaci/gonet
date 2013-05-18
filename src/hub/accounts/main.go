package accounts

import (
	. "types"
	"db/user_tbl"
)

//----------------------------------------------- Add a single user
func AddUser(ud *User) {
	_add_fsm(ud)
	_add_rank(ud)
}

//----------------------------------------------- Load a single usser directly from db
func LoadUser(id int32) error {
	ud, err := user_tbl.Load(id)

	if err == nil {
		_add_fsm(&ud)
		_add_rank(&ud)
	}

	return err
}
