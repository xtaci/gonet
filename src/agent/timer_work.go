package agent

import (
	"fmt"
	"strconv"
	"time"
)

import (
	"cfg"
	. "types"
	"types/estate"
)

//----------------------------------------------- timer work
func timer_work(sess *Session) {
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

	// check whether flush timeout is reached.
	config := cfg.Get()
	ivl, _ := strconv.Atoi(config["flush_interval"])
	if time.Now().Unix()-sess.LastFlush > int64(ivl) {
		fmt.Println("TODO: flush all to db")
		sess.LastFlush = time.Now().Unix()
	}
}
