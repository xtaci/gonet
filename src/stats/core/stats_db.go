package core

import (
	"labix.org/v2/mgo"
)

import (
	"cfg"
)

var _session *mgo.Session

const (
	STATS_COLLECTION = "STATS"
)

func init() {
	config := cfg.Get()
	session, err := mgo.Dial(config["mongo_host_stats"])

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	_session = session
}

func StatsCollection() *mgo.Collection {
	config := cfg.Get()
	return _session.DB(config["mongo_db_stats"]).C(STATS_COLLECTION)
}
