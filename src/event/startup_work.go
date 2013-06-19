package main

import (
//"strconv"
)

import (
//"db/data_tbl"
//"event/core"
//"types/estates"
)

func startup_work() {
	load_estates()
}

//---------------------------------------------------------- load user table into memory
func load_estates() {
	/*
		var managers []estates.Manager
		data_tbl.GetAll(estates.COLLECTION, &managers)

		for _, m := range managers {
			userid := m.UserId
			for event_id, cd := range m.CDs {
				eid, _ := strconv.Atoi(event_id)
				core.Load(estates.COLLECTION, cd.OID, userid, cd.Timeout, uint32(eid))
			}
		}
	*/
}
