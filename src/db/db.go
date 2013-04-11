package db

import "github.com/ziutek/mymysql/mysql"
import _ "github.com/ziutek/mymysql/native" // Native engine

var dbch chan mysql.Conn

func Login(out chan string, name string, password string) {
	db := <-dbch
	stmt := "select id, name, password from users where name = '%s' AND password = MD5('%s')"
	rows,_, err := db.Query(stmt, name, password)
	if err  != nil {
		panic(err.Error())
	}

	if len(rows) >  0 {
		out <- "true"
	} else  {
		out <- "false"
	}

	dbch <- db
}

func Start(max int)  {
	user := "root"
	pass := "qwer1234"
	dbname := "game"

	dbch = make(chan mysql.Conn, max)

	for i:=0;i<max;i++ {
		db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)
		err := db.Connect()

		if err!= nil {
			panic(err)
		}

		dbch <-db
	}
}
