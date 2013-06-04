package core

import (
	"db/user_tbl"
	. "types"
)

//----------------------------------------------- Add a single user
func AddUser(user *User) {
	_add_fsm(user)
	_add_rank(user)
}

//----------------------------------------------- Load a single usser directly from db
func LoadUser(id int32) bool {
	user := user_tbl.Get(id)

	if user != nil {
		_add_fsm(user)
		_add_rank(user)
		return true
	}

	return false
}
