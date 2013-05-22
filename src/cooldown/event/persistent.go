package event

import (
	"db/estate_tbl"
)

func Execute(event *Event) {
	manager := estate_tbl.Get(event.user_id)
	// TODO :perform changes & save back
}
