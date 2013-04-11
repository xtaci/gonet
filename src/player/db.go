package player

import "fmt"
import "strings"
import "reflect"
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
		ud.id = rows[0].Int(0)
		ud.name = rows[0].Str(1)
		out <- "true"
	} else  {
		out <- "false"
	}
}

func (conn * DBConn) Flush(ud *User) {
	v := reflect.ValueOf(ud).Elem()
	key := v.Type()
	count := key.NumField()

	fields := make([]string, count)
	values := make([]string, count)

	slice_idx := 0
	for i := 0; i < count; i++ {
		typeok := true
		switch v.Field(i).Kind() {
		case reflect.Int:
			values[slice_idx] = fmt.Sprintf("'%d'",v.Field(i).Int())
		case reflect.String:
			values[slice_idx] = fmt.Sprintf("'%s'",v.Field(i).String())
		default:
			typeok = false
		}

		if (typeok) {
			fields[slice_idx] = key.Field(i).Name
			slice_idx++
		}
	}

	stmt := []string{"REPLACE INTO cities(", strings.Join(fields[:slice_idx],","),
			 ") VALUES (", strings.Join(values[:slice_idx],","), ")"}

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
