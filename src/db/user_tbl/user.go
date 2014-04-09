package user_tbl

import (
	"crypto/md5"
	"io"
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
	COLLECTION   = "USERS"
	COUNTER_NAME = "USERID_GEN"
)

//---------------------------------------------------------- update a user
func Set(user *User) bool {
	ms, c := C(COLLECTION)
	defer ms.Close()

	info, err := c.Upsert(bson.M{"id": user.Id}, user)
	if err != nil {
		log.Println(info, err)
		return false
	}

	return true
}

//---------------------------------------------------------- login with name & mac address
func LoginMac(name, mac string) *User {
	ms, c := C(COLLECTION)
	defer ms.Close()

	user := &User{}
	err := c.Find(bson.M{"name": name, "mac": mac}).One(user)
	if err != nil {
		log.Println(err, mac)
		return nil
	}

	return user
}

//---------------------------------------------------------- create a new user
func New(name, mac string) *User {
	ms, c := C(COLLECTION)
	defer ms.Close()

	config := cfg.Get()
	user := &User{}
	err := c.Find(bson.M{"name": name}).One(user)
	if err != nil {
		user.Id = NextVal(COUNTER_NAME)
		user.Name = name
		user.Mac = mac
		user.Domain = config["domain"]
		user.CreatedAt = time.Now().Unix()
		err := c.Insert(user)
		if err != nil {
			log.Println(err, name, mac)
			return nil
		}
		return user
	}

	return nil
}

//---------------------------------------------------------- query a user by name
func Query(name string) *User {
	ms, c := C(COLLECTION)
	defer ms.Close()

	user := &User{}
	err := c.Find(bson.M{"name": name}).One(user)
	if err != nil {
		log.Println(err, name)
		return nil
	}

	return user
}

//---------------------------------------------------------- load a user
func Get(id int32) *User {
	ms, c := C(COLLECTION)
	defer ms.Close()

	user := &User{}
	err := c.Find(bson.M{"id": id}).One(user)
	if err != nil {
		log.Println(err, id)
		return nil
	}

	return user
}

//---------------------------------------------------------- load all userss
func GetAll() []User {
	ms, c := C(COLLECTION)
	defer ms.Close()

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
