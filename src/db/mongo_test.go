package db

import (
	"cfg"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"math"
	"testing"
)

type Person struct {
	Name  string
	Phone string
}

func TestMongo(t *testing.T) {
	config := cfg.Get()
	session, err := mgo.Dial(config["mongo_host"])
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		panic(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println("Phone:", result.Phone)
}

func TestOverflow(t *testing.T) {
	config := cfg.Get()
	c := Mongo.DB(config["mongo_db"]).C(COUNTERS)

	counter := &Counter{Name: "testcounter", NextVal: math.MaxInt32}
	c.Remove(bson.M{"name": "testcounter"})
	c.Insert(counter)

	v := NextVal("testcounter")
	fmt.Println(v)

	if v != 0 {
		t.Fatal("failed to overflow int32 to 0")
	}

	v = NextVal("testcounter")
	fmt.Println(v)

	if v != 1 {
		t.Fatal("failed to overflow int32 to 1")
	}
}
