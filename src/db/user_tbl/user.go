package user_tbl

import (
	"crypto/md5"
	"io"
	"labix.org/v2/mgo/bson"
	"log"
)

import (
	"cfg"
	. "db"
	. "types"
)

const (
	COLLECTION   = "USERS"
	COUNTER_NAME = "USERID_GEN"
)

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

func LoginMac(name, mac string) *User {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	user := &User{}
	err := c.Find(bson.M{"name": name, "mac": mac}).One(user)
	if err != nil {
		log.Println(err, mac)
		return nil
	}

	return user
}

//----------------------------------------------- Create a new user
func New(name, mac string) *User {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COLLECTION)

	user := &User{}
	err := c.Find(bson.M{"name": name}).One(user)
	if err != nil {
		user.Id = NextVal(COUNTER_NAME)
		user.Name = name
		user.Mac = mac
		err := c.Insert(user)
		if err != nil {
			log.Println(err, name, mac)
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
		log.Println(err, id)
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
	config := cfg.Get()
	salted := str + config["salt"]
	h := md5.New()
	io.WriteString(h, salted)
	return h.Sum(nil)
}
