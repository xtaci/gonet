package estate_tbl

import (
	"encoding/json"
	"fmt"
	"log"
	"types/estate"
)

import (
	. "db"
)

const (
	PAT_ESTATE = "estate:%v"
)

func Set(user_id int32, manager *estate.Manager) bool {
	json_var, err := json.Marshal(manager)

	if err != nil {
		log.Println(err)
		return false
	}

	set := Redis.Set(fmt.Sprintf(PAT_ESTATE, user_id), string(json_var))
	if set.Err() != nil {
		log.Println(set.Err())
		return false
	}

	return true
}

func Get(user_id int32) *estate.Manager {
	get := Redis.Get(fmt.Sprintf(PAT_ESTATE, user_id))

	if get.Err() != nil {
		log.Println(get.Err())
		return nil
	}

	manager := &estate.Manager{}
	err := json.Unmarshal([]byte(get.Val()), manager)
	if err != nil {
		log.Println(err)
		return nil
	}

	return manager
}
