package res_tbl

import (
	"encoding/json"
	"fmt"
	"log"
	. "types"
)

import (
	. "db"
)

const (
	PAT_RES = "res:%v"
)

func Set(user_id int32, res *Res) bool {
	json_val := res.JSON()

	set := Redis.Set(fmt.Sprintf(PAT_RES, user_id), string(json_val))
	if set.Err() != nil {
		log.Println(set.Err())
		return false
	}

	return true
}

func Get(user_id int32) *Res {
	get := Redis.Get(fmt.Sprintf(PAT_RES, user_id))

	if get.Err() != nil {
		log.Println(get.Err())
		return nil
	}

	res := &Res{}
	err := json.Unmarshal([]byte(get.Val()), res)
	if err != nil {
		log.Println(err)
		return nil
	}

	return res
}
