package event

import (
	. "types"
	playerdata "db/playerdata_tbl"
)

func Execute(event *Event) {
	data := &PlayerData{}
	playerdata.Get(event.user_id, data)
	// TODO :perform changes & save back
}
