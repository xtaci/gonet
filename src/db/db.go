package db

import . "types"
import "strings"
import "strconv"
import "github.com/ziutek/mymysql/mysql"
import _ "github.com/ziutek/mymysql/native" // Native engine
import "log"

const (
	DEFAULT_INSTANCE = 4
)

type DBConn struct {
	dbch chan mysql.Conn
}

func (conn *DBConn) Login(out chan string, name string, password string, ud *User) {
	stmt := "select id, name, password from users where name = '%s' AND password = MD5('%s')"

	db := <-conn.dbch
	rows, res, err := db.Query(stmt, name, password)
	conn.dbch <- db

	if err != nil {
		log.Println(err.Error())
	}

	if len(rows) > 0 {
		sql_load(ud, &rows[0], res)
		out <- "true"
	} else {
		out <- "false"
	}
}

func (conn *DBConn) Flush(ud *User) {
	fields, values := sql_dump(ud)
	stmt := []string{"REPLACE INTO cities(", strings.Join(fields, ","),
		") VALUES (", strings.Join(values, ","), ")"}

	db := <-conn.dbch
	_, _, err := db.Query(strings.Join(stmt, " "))
	conn.dbch <- db
	if err != nil {
		log.Println(err.Error())
	}
}

var DB DBConn
func StartDB(config map[string]string) {
	// instance
	num := DEFAULT_INSTANCE
	if config["max_db_conn"] != "" {
		num,_ = strconv.Atoi(config["max_db_conn"])
	}

	DB.dbch = make(chan mysql.Conn, num)
	log.Println("DB instance:", num)

	for i := 0; i < num; i++ {
		db := mysql.New("tcp", "", config["mysql_host"], config["mysql_username"],
						config["mysql_password"], config["mysql_dbname"])
		err := db.Connect()

		if err != nil {
			log.Panic(err)
		}

		DB.dbch <- db
	}
}
