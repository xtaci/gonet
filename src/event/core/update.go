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
	"types/estates"
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

func _do(event *Event) {
	switch event.tblname {
	case "ESTATES":
		_do_estates(event)
	case "ARMY{":
		_do_army(event)
	}
}

//------------------------------------------------ do the real work until complete or panic!!!
func _do_estates(event *Event) {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C("ESTATES")

	manager := &estates.Manager{}
	err := c.Find(bson.M{"id": event.user_id}).One(manager)

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"version": 1}},
		ReturnNew: true,
	}

	fmt.Println("TODO : change value here")

	// find & update
	_, err = c.Find(bson.M{"id": event.user_id, "version": manager.Version}).Apply(change, manager)
	if err != nil { // repeat
		_do_estates(event)
	}
}

//------------------------------------------------ do the real work until complete or panic!!!
func _do_army(event *Event) {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C("ARMY")

	manager := &estates.Manager{}
	err := c.Find(bson.M{"id": event.user_id}).One(manager)

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"version": 1}},
		ReturnNew: true,
	}

	fmt.Println("TODO : change value here")

	// find & update
	_, err = c.Find(bson.M{"id": event.user_id, "version": manager.Version}).Apply(change, manager)
	if err != nil { // repeat
		_do_army(event)
	}
}
