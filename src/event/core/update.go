package core

import (
	. "helper"
	. "types"
)

//---------------------------------------------------------- perform changes & save back, atomic
func Execute(event *Event, event_id int32) (ret bool) {
	defer PrintPanicStack()
	_do(event, event_id)
	return true
}

func _do(event *Event, event_id int32) {
}
