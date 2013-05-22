package user_tbl

import (
	"fmt"
	"encoding/json"
	"log"
	"crypto/md5"
	"io"
	"time"
)

import (
	. "types"
	. "db"
)

const (
	NEXTVAL = "NEXTVAL"
	PAT_UID = "uid:%v:%v"
	PAT_NAME = "name:%v:uid"
	PAT_BASIC = "basic:%v"
	PAT_BASIC_KEYS = "basic:*"
)


//----------------------------------------------- Change Name in both PAT_UID & PAT_NAME
func ChangeName(basic *Basic, newname string) bool {
	// update uid:1001:name -> xtaci
	set := Redis.Set(fmt.Sprintf(PAT_UID, basic.Id,"name"), newname)
	if set.Err() != nil {
		log.Println(set.Err())
		return false
	}

	// delete name:oldname:uid -> 1001
	del := Redis.Del(fmt.Sprintf(PAT_NAME, basic.Name))
	if del.Err() != nil {
		log.Println(del.Err())
		return false
	}

	// set name:newname:uid -> 1001
	set = Redis.Set(fmt.Sprintf(PAT_NAME,newname), fmt.Sprint(basic.Id))
	if set.Err() != nil {
		log.Println(set.Err())
		return false
	}

	// make changes to basic
	set = Redis.Set(fmt.Sprintf(PAT_BASIC,basic.Id), basic.JSON())
	if set.Err() !=nil {
		log.Println(set.Err())
		return false
	}

	return true
}

//----------------------------------------------- login with (name, password) pair
func Login(name string, pass string) (*Basic) {
	name_uid_key := fmt.Sprintf(PAT_NAME, name)
	_id := Redis.Get(name_uid_key)
	_pass := Redis.Get(fmt.Sprintf(PAT_UID, _id.Val(), "pass"))

	// compare pass
	h := md5.New()
	io.WriteString(h, pass)
	if string(h.Sum(nil)) == _pass.Val() {
		return nil
	}

	basic_json := Redis.Get(fmt.Sprintf(PAT_BASIC,_id.Val()))
	var basic *Basic
	json.Unmarshal([]byte(basic_json.Val()), basic)
	return basic
}

//----------------------------------------------- Create a new user
func New(user, pass string) *Basic {
	// check existence of a user
	exists := Redis.Exists(fmt.Sprintf(PAT_NAME,user))
	if exists.Val() == true {
		log.Println("user exists")
		return nil
	}

	// generate new user id
	next_id := Redis.Incr(NEXTVAL)
	if next_id.Err() !=nil {
		return nil
	}

	id := int32(next_id.Val())

	// uid:1001:name -> xtaci
	set := Redis.Set(fmt.Sprintf(PAT_UID,id,"name"), user)
	if set.Err() != nil {
		log.Println(set.Err())
		return nil
	}

	// uid:1001:pass -> MD5("password")
	h := md5.New()
	io.WriteString(h, pass)
	set = Redis.Set(fmt.Sprintf(PAT_UID,id, "pass"), string(h.Sum(nil)))
	if set.Err() != nil {
		log.Println(set.Err())
		return nil
	}

	// name:xtaci:uid -> 1001
	set = Redis.Set(fmt.Sprintf(PAT_NAME,user), fmt.Sprint(id))
	if set.Err() != nil {
		log.Println(set.Err())
		return nil
	}

	basic := &Basic{}
	basic.Id = id
	basic.Name = user
	basic.LoginCount = 1
	basic.LastLogin = time.Now().Unix()

	set = Redis.Set(fmt.Sprintf(PAT_BASIC,id), basic.JSON())
	if set.Err() !=nil {
		log.Println(set.Err())
		return nil
	}

	return basic
}

//----------------------------------------------- Load a user's basic info 
func Get(id int32) *Basic {
	get := Redis.Get(fmt.Sprintf(PAT_BASIC,id))
	if get.Err() !=nil {
		return nil
	}

	basic := &Basic{}
	err := json.Unmarshal([]byte(get.Val()), basic)
	if err != nil {
		log.Println(err)
		return nil
	}

	return basic
}

//----------------------------------------------- Load all users's basic info
func GetAll() []*Basic {
	get := Redis.Keys(PAT_BASIC_KEYS)
	if get.Err() !=nil {
		return nil
	}

	keys := get.Val()
	basis := make([]*Basic, len(keys))

	for i:=0;i<len(keys);i++ {
		json_val := Redis.Get(keys[i])
		if json_val.Err() !=nil {
			return nil
		}

		basic := &Basic{}
		err := json.Unmarshal([]byte(json_val.Val()), basic)
		if err != nil {
			log.Println(err)
			return nil
		}

		basis[i] = basic
	}

	return basis
}
