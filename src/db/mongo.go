package db

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

import (
	"cfg"
)

var Mongo *mgo.Session

const (
	COUNTERS = "COUNTERS"
)

type Counter struct {
	Name    string
	NextVal int32
}

func init() {
	config := cfg.Get()
	session, err := mgo.Dial(config["mongo_host"])

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	Mongo = session
}

//------------------------------------------------ for very simple use
func Collection(name string) *mgo.Collection {
	config := cfg.Get()
	return Mongo.DB(config["mongo_db"]).C(name)
}

func NextVal(countername string) int32 {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COUNTERS)

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"nextval": 1}},
		Upsert:    true,
		ReturnNew: true,
	}

	next := &Counter{}
	info, err := c.Find(bson.M{"name": countername}).Apply(change, &next)
	if err != nil {
		log.Println(info, err)
		return -1
	}

	return next.NextVal
}
