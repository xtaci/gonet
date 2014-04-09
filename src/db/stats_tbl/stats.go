package stats_tbl

import (
	"cfg"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"time"
)

import (
	. "helper"
	. "types/stats"
)

var _stats_db *mgo.Session // statsdb session

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

//------------------------------------------------------------ 得到玩家日志
func GetPlayLog(userid int, key string, start, end int) []IntGameInfo {
	ms, c := C(INT_GAME_INFO)
	defer ms.Close()
	outs := make([]IntGameInfo, 0)
	k := fmt.Sprintf("%v#%v", userid, key)
	start_time := time.Unix(int64(start), 0).UTC()
	end_time := time.Unix(int64(end), 0).UTC()
	err := c.Find(bson.M{"$and": []bson.M{bson.M{"time": bson.M{"$gte": start_time}}, bson.M{"time": bson.M{"$lt": end_time}}, bson.M{"key": k}}}).All(&outs)
	if err != nil {
		ERR(INT_GAME_INFO, "GetPlayLog", err, userid, key, start, end)
		return nil
	}
	return outs
}

//--------------------------------------------------------- 得到gacha日志
func GetGachaLog(userid, start, end int) []IntGameInfo {
	ms, c := C(INT_GAME_INFO)
	defer ms.Close()
	outs := make([]IntGameInfo, 0)
	k := fmt.Sprintf("^%v#gacha#", userid)
	start_time := time.Unix(int64(start), 0).UTC()
	end_time := time.Unix(int64(end), 0).UTC()
	err := c.Find(bson.M{"$and": []bson.M{bson.M{"time": bson.M{"$gte": start_time}}, bson.M{"time": bson.M{"$lt": end_time}}, bson.M{"key": bson.M{"$regex": k, "$options": "i"}}}}).All(&outs)
	if err != nil {
		ERR(INT_GAME_INFO, "GetGachaLog", err, userid, start, end)
		return nil
	}
	return outs
}

//--------------------------------------------------------- 得到某段时间的adds信息
func GetAdds(start, end int64) []IntGameInfo {
	ms, c := C(INT_GAME_INFO)
	defer ms.Close()
	outs := make([]IntGameInfo, 0)
	err := c.Find(bson.M{"$and": []bson.M{bson.M{"time": bson.M{"$gte": time.Unix(start, 0)}}, bson.M{"time": bson.M{"$lt": time.Unix(end, 0)}}}}).All(&outs)
	if err != nil {
		ERR(INT_GAME_INFO, "GetAdds", err, start, end)
		return nil
	}
	return outs
}

//--------------------------------------------------------- 得到某段时间update信息
func GetStatus(start, end int64) []StrGameInfo {
	ms, c := C(STR_GAME_INFO)
	defer ms.Close()
	outs := make([]StrGameInfo, 0)
	err := c.Find(bson.M{"$and": []bson.M{bson.M{"time": bson.M{"$gte": time.Unix(start, 0)}}, bson.M{"time": bson.M{"$lt": time.Unix(end, 0)}}}}).All(&outs)
	if err != nil {
		ERR(STR_GAME_INFO, "GetStatus", err, start, end)
		return nil
	}
	return outs
}
