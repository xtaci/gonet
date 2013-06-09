package forward_tbl

import (
	"labix.org/v2/mgo/bson"
	"log"
)

import (
	"cfg"
	. "db"
)

const (
	COLLECTION = "FORWARDS"
)

type Forward struct {
	DestId int32
	Value  []byte
}

func Push(dest_id int32, value []byte) bool {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	forward := &Forward{DestId: dest_id, Value: value}
	err := c.Insert(forward)
	if err != nil {
		log.Println(err, dest_id)
		return false
	}

	return true
}

func PopAll(dest_id int32) [][]byte {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	var forwards []Forward
	err := c.Find(bson.M{"destid": dest_id}).All(&forwards)
	if err != nil {
		log.Println(err, dest_id)
	}

	info, err := c.RemoveAll(bson.M{"destid": dest_id})
	if err != nil {
		log.Println(info, err, dest_id)
	}

	objs := make([][]byte, len(forwards))
	for k := range forwards {
		objs[k] = forwards[k].Value
	}

	return objs
}
