package core

import (
	"db/user_tbl"
	. "types"
)

//----------------------------------------------- load all users into memory
func LoadAllUsers() {
	// load users
	uds := user_tbl.GetAll()

	for i := range uds {
		_add_user(&uds[i])
	}
}

//----------------------------------------------- load a single user
func LoadUser(id int32) bool {
	user := user_tbl.Get(id)

	if user != nil {
		_add_user(user)
		return true
	}

	return false
}

//----------------------------------------------- feed user to data structures
func _add_user(user *User) {
	_add_fsm(user)
}
