package accounts

import (
	. "types"
)

func AddUser(ud *User) {
	_add_fsm(ud)
	_add_rank(ud)
}
