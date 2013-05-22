package db

import (
	"strconv"
	"cfg"
)

import (
	"github.com/vmihailenco/redis"
)

var Redis *redis.Client

func init() {
	config := cfg.Get()
	db := -1
	db,_ = strconv.Atoi(config["redis_db"])
	Redis = redis.NewTCPClient(config["redis_host"], config["redis_pass"], int64(db))
}

func CAS(multi *redis.MultiClient, key, value string) ([]redis.Req, error) {
    get := multi.Get(key)
    if err := get.Err(); err != nil && err != redis.Nil {
        return nil, err
    }

    reqs, err := multi.Exec(func() {
        multi.Set(key, value)
    })

    // Transaction failed. Repeat.
    if err == redis.Nil {
        return CAS(multi, key, value)
    }
    return reqs, err
}
