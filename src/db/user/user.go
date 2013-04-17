package user

import . "db"
import . "types"
import "strings"
import "fmt"
import "log"

func Flush(ud *User) {
	fields, values := SQL_dump(ud)
	changes := SQL_set_clause(fields,values)

	stmt := []string{"UPDATE users SET ", strings.Join(changes, ","), " WHERE id=", fmt.Sprint(ud.Id)}

	db := <-DBCH
	defer func(){DBCH <- db}()
	_, _, err := db.Query(strings.Join(stmt, " "))
	if err != nil {
		log.Println(err.Error())
	}
}

func Login(out chan string, name string, password string, ud *User) {
	stmt := "select * from users where name = '%s' AND password = MD5('%s')"

	db := <-DBCH
	defer func(){DBCH <- db}()
	rows, res, err := db.Query(stmt, SQL_escape(name), SQL_escape(password))

	if err != nil {
		log.Println(err.Error())
	}

	if len(rows) > 0 {
		SQL_load(ud, &rows[0], res)
		out <- "true"
		// fake cities
		ud.Cities = make([]City, 1)
		ud.Cities[0].Name = "city of" + ud.Name
		ud.Cities[0].OwnerId = ud.Id
	} else {
		out <- "false"
	}
}
