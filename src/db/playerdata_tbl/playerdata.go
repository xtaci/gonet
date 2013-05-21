package playerdata_tbl

import (
	"github.com/hoisie/redis"
	"cfg"
	. "types"
	"encoding/json"
	"log"
	"fmt"
)

var _redis redis.Client

func init() {
	config := cfg.Get()
	_redis.Addr = config["redis_host"]
}

func Set(user_id int32, data *PlayerData) bool {
	json_var, err := json.Marshal(data)

	if err!= nil {
		log.Println(err)
		return false
	}

	err = _redis.Set(fmt.Sprintf("DATA#%v",user_id), json_var)
	if err!=nil {
		log.Println(err)
		return false
	}

	return true
}

func Get(user_id int32, data *PlayerData) bool {
	json_val, err := _redis.Get(fmt.Sprintf("DATA#%v",user_id))

	if err !=nil {
		log.Println(err)
		return false
	}

	err = json.Unmarshal(json_val, data)
	if err !=nil {
		log.Println(err)
		return false
	}

	return true
}
