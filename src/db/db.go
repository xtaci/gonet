package db

import (
	"log"
	"strconv"
)

import (
	"cfg"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
)

const (
	DEFAULT_INSTANCE = 4
)

var DBCH chan mysql.Conn

func StartDB() {

	config := cfg.Get()
	// instance
	num := DEFAULT_INSTANCE
	if config["max_db_conn"] != "" {
		num, _ = strconv.Atoi(config["max_db_conn"])
	}

	DBCH = make(chan mysql.Conn, num)
	log.Println("DB instance:", num)

	for i := 0; i < num; i++ {
		db := mysql.New("tcp", "", config["mysql_host"], config["mysql_username"],
			config["mysql_password"], config["mysql_dbname"])
		err := db.Connect()

		CheckErr(err)

		DBCH <- db
	}
}
