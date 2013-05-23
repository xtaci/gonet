package event

import (
	"fmt"
	"log"
)

import (
	"db"
	"db/estate_tbl"
	"github.com/vmihailenco/redis"
)

//------------------------------------------------ perform changes & save back, atomic
func Execute(event *Event) {
	multi, _ := db.Redis.MultiClient()
	defer multi.Close()

	key := fmt.Sprintf(estate_tbl.PAT_ESTATE, event.user_id)
	watch := multi.Watch(key)
	_ = watch.Err()

	reqs, err := _do(multi, key)
	if err!= nil {
		log.Println(err,reqs)
	}
}

func _do(multi *redis.MultiClient, key string) ([]redis.Req, error) {
	get := multi.Get(key)
	if err := get.Err(); err != nil && err != redis.Nil {
		return nil, err
	}

	reqs, err := multi.Exec(func() {
		fmt.Println("TODO : change value here")
		value := get.Val()
		multi.Set(key, value)
	})

	// Transaction failed. Repeat.
	if err == redis.Nil {
		return _do(multi, key)
	}
	return reqs, err
}
