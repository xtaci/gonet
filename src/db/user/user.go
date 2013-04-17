package user

import . "db"
import . "types"
import "strings"
import "fmt"

func Flush(ud *User) {
	fields, values := SQL_dump(ud)
	changes := SQL_set_clause(fields,values)

	stmt := []string{"UPDATE users SET ", strings.Join(changes, ","), " WHERE id=", fmt.Sprint(ud.Id)}

	db := <-DBCH
	defer func(){DBCH <- db}()
	_, _, err := db.Query(strings.Join(stmt, " "))

	CheckErr(err)
}

func Login(out chan string, name string, password string, ud *User) {
	stmt := "select * from users where name = '%v' AND password = MD5('%v')"

	db := <-DBCH
	defer func(){DBCH <- db}()
	rows, res, err := db.Query(stmt, SQL_escape(name), SQL_escape(password))

	CheckErr(err)

	if len(rows) > 0 {
		SQL_load(ud, &rows[0], res)
		out <- "true"
	} else {
		out <- "false"
	}
}
