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
