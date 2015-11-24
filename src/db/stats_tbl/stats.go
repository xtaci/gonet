package stats_tbl

import (
	"cfg"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"time"
)

import (
	. "helper"
)

var _stats_db *mgo.Session // statsdb session

const (
	INT_GAME_INFO = "INT_GAME_INFO"
	STR_GAME_INFO = "STR_GAME_INFO"
)

type IntGameInfo struct {
	IntValue int32
	Key      string
	Time     time.Time
	Lang     string
}

type StrGameInfo struct {
	StrValue string
	Key      string
	Time     time.Time
	Lang     string
}

func init() {
	config := cfg.Get()
	// dial mongodb
	sess, err := mgo.Dial(config["mongo_host_stats"])
	if err != nil {
		ERR(err)
		os.Exit(-1)
	}

	// set default session mode to eventual
	sess.SetMode(mgo.Eventual, true)
	_stats_db = sess
}

//------------------------------------------------ copy connection
func C(collection string) (*mgo.Session, *mgo.Collection) {
	config := cfg.Get()
	ms := _stats_db.Copy()
	c := ms.DB(config["mongo_db_stats"]).C(collection)
	return ms, c
}

//------------------------------------------------------------ 记录玩家累计信息
func SetAdds(key string, value int32, lang string) bool {
	add_info := IntGameInfo{
		IntValue: value,
		Key:      key,
		Time:     time.Now().UTC(),
		Lang:     lang,
	}
	ms, c := C(INT_GAME_INFO)
	defer ms.Close()

	err := c.Insert(add_info)
	if err != nil {
		ERR(INT_GAME_INFO, "SetAdds", err, key, value, lang)
		return false
	}
	return true
}

//------------------------------------------------------------ 记录玩家状态信息
func SetUpdate(key, value, lang string) bool {
	ms, c := C(STR_GAME_INFO)
	defer ms.Close()
	info := StrGameInfo{value, key, time.Now().UTC(), lang}
	_, err := c.Upsert(bson.M{"key": key}, info)
	if err != nil {
		ERR(STR_GAME_INFO, "SetUpdate", err, key, value, lang)
	}
	return true
}
