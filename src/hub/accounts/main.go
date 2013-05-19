package accounts

import (
	"db/user_tbl"
	. "types"
)

//----------------------------------------------- Add a single user
func AddUser(ud *User) {
	_add_fsm(ud)
	_add_rank(ud)
}

//----------------------------------------------- Load a single usser directly from db
func LoadUser(id int32) bool {
	ud, err := user_tbl.Load(id)

	if err == nil {
		_add_fsm(&ud)
		_add_rank(&ud)
		return true
	}

	return false
}
