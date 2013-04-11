package db

import "sync"
import "time"
import "database/sql"
import "container/list"
import _ "github.com/go-sql-driver/mysql"

var db_list *list.List
var db_lock sync.Mutex

func _get_db() *sql.DB {
	for {
		db_lock.Lock()
		db := db_list.Front()
		if db != nil {
			db_list.Remove(db)
			db_lock.Unlock()
			return db.Value.(*sql.DB)
		} else {
			db_lock.Unlock()
			time.Sleep(time.Millisecond)
		}
	}

	return nil
}

func _push_db(db *sql.DB)  {
	db_lock.Lock()
	db_list.PushBack(db)
	db_lock.Unlock()
}

func Login(out chan string, name string, password string) {
	db := _get_db()
	stmt := "select id, name, password from users where name = ? AND password = MD5(?)"
	rows, err := db.Query(stmt, name, password)
	if err  != nil {
		panic(err.Error())
	}

	if rows.Next() {
		out <- "true"
	} else  {
		out <- "false"
	}

	_push_db(db)
}

func Start(num_inst int)  {
	db_list = list.New()

	for i:=0;i<num_inst;i++ {
		db, err := sql.Open("mysql", "root:qwer1234@/game?charset=latin1")
		if  err != nil  {
			print("Cannoe connect to db %s", err)
			return
		}
		db_list.PushFront(db)
	}
}
