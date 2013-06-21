package forward_tbl

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

import (
	"cfg"
	. "db"
	. "types"
)

const (
	COLLECTION = "FORWARDS"
)

//---------------------------------------------------------- push an ipc object to db
func Push(req *IPCObject) bool {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	err := c.Insert(req)
	if err != nil {
		log.Println(err, req)
		return false
	}

	return true
}

//---------------------------------------------------------- pop all message for dest user
func PopAll(dest_id int32) []IPCObject {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	var objects []IPCObject
	change := mgo.Change{
		Update: bson.M{"$set": bson.M{"markdelete": true}},
	}

	// mark delete
	info, err := c.Find(bson.M{"destid": dest_id}).Apply(change, &objects)
	if err != nil {
		log.Println(err, info)
	}

	// select
	err = c.Find(bson.M{"destid": dest_id, "markdelete": true}).Sort("time").All(&objects)
	if err != nil {
		log.Println(err)
	}

	// real delete
	info, err = c.RemoveAll(bson.M{"destid": dest_id, "markdelete": true})
	if err != nil {
		log.Println(err, info)
	}

	return objects
}
