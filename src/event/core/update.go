package core

import (
//"fmt"
//"labix.org/v2/mgo"
//"labix.org/v2/mgo/bson"
//"log"
)

import (
	//"cfg"
	//	. "db"
	//	"misc/naming"
	. "helper"

//	"types/estates"
)

//------------------------------------------------ perform changes & save back, atomic
func Execute(event *Event, event_id uint32) (ret bool) {
	defer PrintPanicStack()
	_do(event, event_id)
	return true
}

func _do(event *Event, event_id uint32) {
	/*
		switch event.tblname {
		case naming.FNV1a(estates.COLLECTION):
			_do_estates(event, event_id)
		}
	*/
}

/*

//------------------------------------------------ do the real work until complete or panic!!!
func _do_estates(event *Event, event_id uint32) {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(estates.COLLECTION)

	manager := &estates.Manager{}
	err := c.Find(bson.M{"userid": event.user_id}).One(manager)
	cur_version := manager.Version
	manager.Version++

	// prerequestiques check
	if manager.CDs == nil {
		return
	}

	if manager.Estates == nil {
		return
	}

	EventId := fmt.Sprint(event_id)
	if CD := manager.CDs[EventId]; CD != nil {
		if E := manager.Estates[fmt.Sprint(CD.OID)]; E != nil {
			E.Status = estates.STATUS_NORMAL
		}
		delete(manager.CDs, EventId)
	}

	// find & update
	change := mgo.Change{
		Update:    manager,
		ReturnNew: true,
	}
	_, err = c.Find(bson.M{"userid": event.user_id, "version": cur_version}).Apply(change, manager)
	if err != nil { // repeat
		_do_estates(event, event_id)
	}

	log.Printf("Event Exec: [%v]\n", event_id)
}
*/
