package res_tbl

import (
	"labix.org/v2/mgo/bson"
	"log"
)

import (
	. "db"
	. "types"
	"cfg"
)

const (
	COLLECTION = "RES"
)

func Set(res *Res) bool {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	res.Version++
	info, err := c.Upsert(bson.M{"id": res.Id}, res)
	if err != nil {
		log.Println(info, err)
		return false
	}

	return true
}

func Get(user_id int32) *Res {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	res := &Res{}
	err := c.Find(bson.M{"id": user_id}).One(res)
	if err != nil {
		log.Println(err)
		return nil
	}

	return res
}
