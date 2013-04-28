package user_tbl

import (
	. "db"
	. "types"
)

import (
	"strings"
	"fmt"
	"errors"
)

func Flush(ud *User) {
	fields, values := SQL_dump(ud)
	changes := SQL_set_clause(fields, values)

	stmt := []string{"UPDATE users SET ", strings.Join(changes, ","), " WHERE id=", fmt.Sprint(ud.Id)}

	db := <-DBCH
	defer func() { DBCH <- db }()
	_, _, err := db.Query(strings.Join(stmt, " "))

	CheckErr(err)
}

func Login(name string, password string, ud *User) bool {
	stmt := "select * from users where name = '%v' AND password = MD5('%v')"

	db := <-DBCH
	defer func() { DBCH <- db }()
	rows, res, err := db.Query(stmt, SQL_escape(name), SQL_escape(password))

	CheckErr(err)

	if len(rows) > 0 {
		SQL_load(ud, &rows[0], res)
		return true
	}

	return false
}

func LoginMAC(mac string, ud *User) bool {
	stmt := "SELECT * FROM users where mac='%v'"

	db := <-DBCH
	defer func() { DBCH <- db }()
	rows, res, err := db.Query(stmt, mac)

	CheckErr(err)

	if len(rows) > 0 {
		SQL_load(ud, &rows[0], res)
		return true
	}

	return false
}

func New(ud *User) (ret bool) {
	defer func() {
		if x := recover(); x != nil {
			ret = false
		}
	}()

	fields, values := SQL_dump(ud, "id")
	stmt := []string{"INSERT INTO users(", strings.Join(fields, ","),
		") VALUES (", strings.Join(values, ","), ")"}

	db := <-DBCH
	defer func() { DBCH <- db }()
	_, res, err := db.Query(strings.Join(stmt, " "))
	CheckErr(err)
	ud.Id = int32(res.InsertId())

	if err == nil {
		return true
	}
	return false
}


func Read(id int32) (ud User, err error) {
	stmt := "SELECT * FROM users where id ='%v'"

	db := <-DBCH
	defer func() { DBCH <- db }()

	rows, res, err := db.Query(stmt, id)
	CheckErr(err)

	if len(rows) > 0 {
		SQL_load(&ud, &rows[0], res)
		return
	}

	err = errors.New(fmt.Sprint("cannot find user with id:%v" , id))
	return
}
