package forward_tbl

import (
	"labix.org/v2/mgo/bson"
	"log"
)

import (
	. "db"
	. "types"
)

const (
	COLLECTION = "FORWARDS"
)

//---------------------------------------------------------- push an ipc object to db
func Push(req *IPCObject) bool {
	ms, c := C(COLLECTION)
	defer ms.Close()
	req.MarkDelete = false
	err := c.Insert(req)
	if err != nil {
		log.Println(err, req)
		return false
	}

	return true
}

//---------------------------------------------------------- pop all message for dest user
func PopAll(dest_id int32) []IPCObject {
	ms, c := C(COLLECTION)
	defer ms.Close()

	var objects []IPCObject
	// mark delete
	info, err := c.UpdateAll(bson.M{"destid": dest_id}, bson.M{"$set": bson.M{"markdelete": true}})
	if err != nil {
		log.Println(err, info)
	}

	// select
	err = c.Find(bson.M{"destid": dest_id, "markdelete": true}).Sort("$natural").All(&objects)
	if err != nil {
		log.Println(err)
	}

	// real delete
	info, err = c.RemoveAll(bson.M{"destid": dest_id, "markdelete": true})
	if err != nil {
		log.Println(err, info)
	}

	return objects
}
