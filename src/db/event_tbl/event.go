package event_tbl

import (
	"labix.org/v2/mgo/bson"
	"log"
)

import (
	"cfg"
	. "db"
)

const (
	COLLECTION = "EVENTS"
)

type Event struct {
	UserId  int32
	EventId int32
	Type    int16
	Params  []byte
}

func Push(event_id, user_id int32, _type int16, params []byte) bool {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	event := &Event{UserId: user_id, EventId: event_id, Type: _type, Params: params}
	err := c.Insert(event)
	if err != nil {
		log.Println(err, user_id)
		return false
	}

	return true
}

//----------------------------------------------------------  remove an event
func Remove(event_id int32) bool {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	err := c.Remove(bson.M{"eventid": event_id})
	if err != nil {
		log.Println(err, event_id)
		return false
	}

	return true
}

//----------------------------------------------------------  get an event
func Get(event_id int32) *Event {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	event := &Event{}
	err := c.Find(bson.M{"eventid": event_id}).One(event)
	if err != nil {
		log.Println(err, event_id)
		return nil
	}

	return event
}

//----------------------------------------------------------  get all events
func GetAll() []Event {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	var events []Event
	err := c.Find(nil).All(&events)
	if err != nil {
		log.Println(err)
	}

	return events
}
