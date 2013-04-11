package player

import "fmt"
import "strings"
import "reflect"
import "github.com/ziutek/mymysql/mysql"
import _ "github.com/ziutek/mymysql/native" // Native engine

var dbch chan mysql.Conn

func DBLogin(out chan string, name string, password string, ud * UserData) {
	stmt := "select id, name, password from users where name = '%s' AND password = MD5('%s')"
	db := <-dbch
	rows,_, err := db.Query(stmt, name, password)
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

	dbch <- db
}

func DBFlush(ud *UserData) {
	key := reflect.TypeOf(ud).Elem()
	v := reflect.ValueOf(ud).Elem()
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

	db := <-dbch
	_,_, err := db.Query(strings.Join(stmt, " "))
	if err  != nil {
		println(err.Error())
	}
	dbch <- db
}

func DBStart(max int)  {
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
