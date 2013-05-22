package accounts

import (
	"db/user_tbl"
	. "types"
)

//----------------------------------------------- Add a single user
func AddUser(basic *Basic) {
	_add_fsm(basic)
	_add_rank(basic)
}

//----------------------------------------------- Load a single usser directly from db
func LoadUser(id int32) bool {
	basic := user_tbl.Get(id)

	if basic != nil {
		_add_fsm(basic)
		_add_rank(basic)
		return true
	}

	return false
}
