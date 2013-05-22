package estate_tbl

import (
	"github.com/vmihailenco/redis"
	"cfg"
	"types/estate"
	"encoding/json"
	"log"
	"fmt"
)

const (
	PAT_DATA = "ESTATE:%v"
)

var _redis *redis.Client

func init() {
	config := cfg.Get()
	_redis = redis.NewTCPClient(config["redis_host"], "", -1)
}

func Set(user_id int32, manager *estate.EstateManager) bool {
	json_var, err := json.Marshal(manager)

	if err!= nil {
		log.Println(err)
		return false
	}

	set := _redis.Set(fmt.Sprintf(PAT_DATA,user_id), string(json_var))
	if set.Err() !=nil {
		log.Println(set.Err())
		return false
	}

	return true
}

func Get(user_id int32) (*estate.EstateManager) {
	get := _redis.Get(fmt.Sprintf(PAT_DATA,user_id))

	if get.Err() !=nil {
		log.Println(get.Err())
		return nil
	}

	manager := &estate.EstateManager{}
	err := json.Unmarshal([]byte(get.Val()), manager)
	if err !=nil {
		log.Println(err)
		return nil
	}

	return manager
}
