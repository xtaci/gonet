package core

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

import (
	"cfg"
	. "db"
	"db/estate_tbl"
	"types/estate"
)

//------------------------------------------------ perform changes & save back, atomic
func Execute(event *Event) (ret bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
			ret = false
		}
	}()

	_do(event)
	return true
}

//------------------------------------------------ do the real work until complete or panic!!!
func _do(event *Event) {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(estate_tbl.COLLECTION)

	manager := &estate.Manager{}
	err := c.Find(bson.M{"id": event.user_id}).One(manager)

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"version": 1}},
		ReturnNew: true,
	}

	fmt.Println("TODO : change value here")

	// find & update
	_, err = c.Find(bson.M{"id": event.user_id, "version": manager.Version}).Apply(change, manager)
	if err != nil { // repeat
		_do(event)
	}
}
