package agent

import (
	"time"
)

import (
	. "types"
	"types/estate"
)

//----------------------------------------------- session timeout
func event_work(sess *Session, session_timeout int) bool {
	// check building upgrades
	CDs := sess.EstateManager.CDs
	Estates := sess.EstateManager.Estates
	expire_cds := make([]bool, len(CDs))
	for i := range CDs {
		if CDs[i].Timeout <= time.Now().Unix() {	// times up
			expire_cds[i] = true
			for k := range Estates {
				if CDs[i].OID == Estates[k].OID {	// if it is the oid
					Estates[k].Status = estate.STATUS_NORMAL
				}
			}
		}
	}

	// update CDs
	var updated []estate.CD
	for k,v :=range expire_cds {
		if !v {
			updated = append(updated, CDs[k])
		}
	}

	sess.EstateManager.CDs = updated

	return false
}
