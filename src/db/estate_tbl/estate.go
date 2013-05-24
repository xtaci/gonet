package estate_tbl

import (
	"log"
	"labix.org/v2/mgo/bson"
)

import (
	. "db"
	"cfg"
	"types/estate"
)

const (
	PAT_ESTATE = "estate:%v"
	COLLECTION = "ESTATE"
)

func Set(manager *estate.Manager) bool {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	manager.Version++
	info, err := c.Upsert(bson.M{"id":manager.Id}, manager)
	if err != nil {
		log.Println(info, err)
		return false
	}

	return true
}

func Get(user_id int32) *estate.Manager {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	manager := &estate.Manager{}
	err := c.Find(bson.M{"id":user_id}).One(manager)
	if err != nil {
		log.Println(err)
		return nil
	}

	return manager
}
