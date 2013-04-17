package db

import . "types"
import "strconv"
import "github.com/ziutek/mymysql/mysql"
import _ "github.com/ziutek/mymysql/native" // Native engine
import "log"

const (
	DEFAULT_INSTANCE = 4
)

var DBCH chan mysql.Conn

func Login(out chan string, name string, password string, ud *User) {
	stmt := "select * from users where name = '%s' AND password = MD5('%s')"

	db := <-DBCH
	defer func(){DBCH <- db}()
	rows, res, err := db.Query(stmt, sql_escape(name), sql_escape(password))

	if err != nil {
		log.Println(err.Error())
	}

	if len(rows) > 0 {
		sql_load(ud, &rows[0], res)
		out <- "true"
		// fake cities
		ud.Cities = make([]City, 1)
		ud.Cities[0].Name = "city of" + ud.Name
		ud.Cities[0].OwnerId = ud.Id
	} else {
		out <- "false"
	}
}



func StartDB(config map[string]string) {
	// instance
	num := DEFAULT_INSTANCE
	if config["max_db_conn"] != "" {
		num,_ = strconv.Atoi(config["max_db_conn"])
	}

	DBCH = make(chan mysql.Conn, num)
	log.Println("DB instance:", num)

	for i := 0; i < num; i++ {
		db := mysql.New("tcp", "", config["mysql_host"], config["mysql_username"],
						config["mysql_password"], config["mysql_dbname"])
		err := db.Connect()

		if err != nil {
			log.Panic(err)
		}

		DBCH <- db
	}
}
