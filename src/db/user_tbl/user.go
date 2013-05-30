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
func Set(user *User) bool {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	info, err := c.Upsert(bson.M{"id": user.Id}, user)
	if err != nil {
		log.Println(info, err)
		return false
	}

	return true
}

//----------------------------------------------- login with (name, password) pair
func Login(name string, pass string) *User {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	user := &User{}
	err := c.Find(bson.M{"name": name, "pass": _md5(pass)}).One(user)
	if err != nil {
		log.Println(err)
		return nil
	}

	return user
}

//----------------------------------------------- Create a new user
func New(name, pass string) *User {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	user := &User{}
	err := c.Find(bson.M{"name": name}).One(user)
	if err != nil {
		user.Id = _nextval()
		user.Name = name
		user.Pass = _md5(pass)
		user.LoginCount = 1
		user.LastLogin = time.Now().Unix()
		err := c.Insert(user)
		if err != nil {
			log.Println(err)
			return nil
		}
		return user
	}

	return nil
}

//----------------------------------------------- Load a user's user info
func Get(id int32) *User {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	user := &User{}
	err := c.Find(bson.M{"id": id}).One(user)
	if err != nil {
		log.Println(err)
		return nil
	}

	return user
}

//----------------------------------------------- Load all users's user info
func GetAll() []User {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	var users []User
	err := c.Find(nil).All(&users)
	if err != nil {
		log.Println(err)
		return nil
	}

	return users
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
