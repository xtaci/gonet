package estate_tbl

import (
	"types/estate"
	"encoding/json"
	"log"
	"fmt"
)

import (
	. "db"
)

const (
	PAT_ESTATE = "estate:%v"
)

func Set(user_id int32, manager *estate.Manager) bool {
	json_var, err := json.Marshal(manager)

	if err!= nil {
		log.Println(err)
		return false
	}

	// CAS operation
	multi, _ := Redis.MultiClient()
	defer multi.Close()

	key :=  fmt.Sprintf(PAT_ESTATE, user_id)
	watch := multi.Watch(key)
	_ = watch.Err()

	reqs, err := CAS(multi, key, string(json_var))
	if err!=nil {
		log.Println(reqs, err)
		return false
	}

	return true
}

func Get(user_id int32) (*estate.Manager) {
	get := Redis.Get(fmt.Sprintf(PAT_ESTATE,user_id))

	if get.Err() !=nil {
		log.Println(get.Err())
		return nil
	}

	manager := &estate.Manager{}
	err := json.Unmarshal([]byte(get.Val()), manager)
	if err !=nil {
		log.Println(err)
		return nil
	}

	return manager
}
