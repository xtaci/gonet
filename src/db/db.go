package db

import . "types"
import "strings"
import "github.com/ziutek/mymysql/mysql"
import _ "github.com/ziutek/mymysql/native" // Native engine


type DBConn struct {
	dbch chan mysql.Conn;
}

func (conn * DBConn) Login(out chan string, name string, password string, ud * User) {
	stmt := "select id, name, password from users where name = '%s' AND password = MD5('%s')"

	db := <-conn.dbch
	rows,_, err := db.Query(stmt, name, password)
	conn.dbch <- db

	if err  != nil {
		panic(err.Error())
	}

	if len(rows) >  0 {
		ud.Id = rows[0].Int(0)
		ud.Name = rows[0].Str(1)
		out <- "true"
	} else  {
		out <- "false"
	}
}

func (conn * DBConn) Flush(ud * User) {
	fields, values := sql_dump(ud)
	stmt := []string{"REPLACE INTO cities(", strings.Join(fields,","),
			 ") VALUES (", strings.Join(values,","), ")"}

	db := <-conn.dbch
	_,_, err := db.Query(strings.Join(stmt, " "))
	conn.dbch <-db
	if err  != nil {
		println(err.Error())
	}
}

var DB DBConn

func StartDB(max int) {
	user := "root"
	pass := "qwer1234"
	dbname := "game"

	DB.dbch = make(chan mysql.Conn, max)

	for i:=0;i<max;i++ {
		db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
		err := db.Connect()

		if err!= nil {
			panic(err)
		}

		DB.dbch <-db
	}
}
