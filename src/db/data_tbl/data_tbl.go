package data_tbl

import (
	"labix.org/v2/mgo/bson"
	"log"
	"reflect"
)

import (
	. "db"
)

//------------------------------------------------ pass-in *ptr
func Set(collection string, data interface{}) bool {
	ms, c := C(collection)
	defer ms.Close()
	v := reflect.ValueOf(data).Elem()

	version := v.FieldByName("Version")

	if !version.IsValid() {
		log.Println(`Cannot seriazlie a struct without "Version" Field`)
		return false
	}
	version.SetUint(uint64(version.Interface().(uint32) + 1))

	id := v.FieldByName("UserId")

	if !id.IsValid() {
		log.Println(`Cannot seriazlie a struct without "UserId" Field`)
		return false
	}

	info, err := c.Upsert(bson.M{"userid": id.Interface().(int32)}, data)
	if err != nil {
		log.Println(info, err)
		return false
	}

	return true
}

//------------------------------------------------ pass-in  *ptr or **ptr
func Get(collection string, user_id int32, data interface{}) bool {
	ms, c := C(collection)
	defer ms.Close()

	err := c.Find(bson.M{"userid": user_id}).One(data)
	if err != nil {
		log.Println(err, collection, user_id)
		return false
	}

	return true
}

//------------------------------------------------ pass-in *[]slice
func GetAll(collection string, all interface{}) bool {
	ms, c := C(collection)
	defer ms.Close()

	err := c.Find(nil).All(all)
	if err != nil {
		log.Println(err, collection)
		return false
	}

	return true
}
