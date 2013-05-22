package playerdata_tbl

import (
	"github.com/vmihailenco/redis"
	"cfg"
	. "types"
	"encoding/json"
	"log"
	"fmt"
)

var _redis *redis.Client

func init() {
	config := cfg.Get()
	_redis = redis.NewTCPClient(config["redis_host"], "", -1)
}

func Set(user_id int32, data *PlayerData) bool {
	json_var, err := json.Marshal(data)

	if err!= nil {
		log.Println(err)
		return false
	}

	set := _redis.Set(fmt.Sprintf("DATA#%v",user_id), string(json_var))
	if set.Err() !=nil {
		log.Println(set.Err())
		return false
	}

	return true
}

func Get(user_id int32, data *PlayerData) bool {
	get := _redis.Get(fmt.Sprintf("DATA#%v",user_id))

	if get.Err() !=nil {
		log.Println(get.Err())
		return false
	}

	err := json.Unmarshal([]byte(get.Val()), data)
	if err !=nil {
		log.Println(err)
		return false
	}

	return true
}
