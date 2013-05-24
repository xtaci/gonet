package db

import (
	"cfg"
)

import (
	"labix.org/v2/mgo"
)

var Mongo mgo.Session

func init() {
	config := cfg.Get()
	Mongo, err := mgo.Dial(config["mongo_host"])

	if err != nil {
		panic(err)
	}

	Mongo.SetMode(mgo.Monotonic, true)
}
