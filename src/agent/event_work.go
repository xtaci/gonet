package agent

import (
	"time"
)

import (
	. "types"
	"types/estate"
)

//----------------------------------------------- session timeout
func event_work(sess *Session) {
	// check building upgrades
	CDs := sess.EstateManager.CDs
	Estates := sess.EstateManager.Estates

	for i := range CDs {
		if CDs[i].Timeout <= time.Now().Unix() { // times up
			for k := range Estates {
				if CDs[i].OID == Estates[k].OID { // if it is the oid
					Estates[k].Status = estate.STATUS_NORMAL
				}
			}
			delete(CDs, i)
		}
	}
}
