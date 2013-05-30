package defensive_tbl

import (
	"labix.org/v2/mgo/bson"
	"log"
)

import (
	"cfg"
	. "db"
	"types/defensive"
)

const (
	COLLECTION = "DEFENSIVE"
)

func Set(manager *defensive.Manager) bool {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	manager.Version++
	info, err := c.Upsert(bson.M{"id": manager.Id}, manager)
	if err != nil {
		log.Println(info, err)
		return false
	}

	return true
}

func Get(user_id int32) *defensive.Manager {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	manager := &defensive.Manager{}
	err := c.Find(bson.M{"id": user_id}).One(manager)
	if err != nil {
		log.Println(err)
		return nil
	}

	return manager
}
