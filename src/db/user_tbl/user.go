package user_tbl

import (
	"crypto/md5"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

import (
	"cfg"
	. "db"
	. "types"
)

const (
	COLLECTION = "BASIC"
	NEXTVAL    = "NEXTVAL"
)

type NextVal struct {
	ID int32
}

//----------------------------------------------- Change
func Set(basic *Basic) bool {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	info, err := c.Upsert(bson.M{"id": basic.Id}, basic)
	if err != nil {
		log.Println(info, err)
		return false
	}

	return true
}

//----------------------------------------------- login with (name, password) pair
func Login(name string, pass string) *Basic {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	basic := &Basic{}
	err := c.Find(bson.M{"name": name, "pass": _md5(pass)}).One(basic)
	if err != nil {
		log.Println(err)
		return nil
	}

	return basic
}

//----------------------------------------------- Create a new user
func New(name, pass string) *Basic {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	basic := &Basic{}
	err := c.Find(bson.M{"name": name}).One(basic)
	if err != nil {
		basic.Id = _nextval()
		basic.Name = name
		basic.Pass = _md5(pass)
		basic.LoginCount = 1
		basic.LastLogin = time.Now().Unix()
		err := c.Insert(basic)
		if err != nil {
			log.Println(err)
			return nil
		}
		return basic
	}

	return nil
}

//----------------------------------------------- Load a user's basic info
func Get(id int32) *Basic {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	basic := &Basic{}
	err := c.Find(bson.M{"id": id}).One(basic)
	if err != nil {
		log.Println(err)
		return nil
	}

	return basic
}

//----------------------------------------------- Load all users's basic info
func GetAll() []Basic {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	var basis []Basic
	err := c.Find(nil).All(&basis)
	if err != nil {
		log.Println(err)
		return nil
	}

	return basis
}

func _md5(str string) []byte {
	h := md5.New()
	io.WriteString(h, str)
	return h.Sum(nil)
}

func _nextval() int32 {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(NEXTVAL)

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"n": 1}},
		ReturnNew: true,
	}

	next := &NextVal{}
	info, err := c.Find(nil).Apply(change, &next)
	if err != nil {
		log.Println(info, err)
		return -1
	}

	return next.ID
}
