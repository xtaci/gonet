package db

import . "types"
import "strings"
import "fmt"
import "log"

func UserFlush(ud *User) {
	fields, values := sql_dump(ud)
	changes := to_set_clause(fields,values)

	stmt := []string{"UPDATE users SET ", strings.Join(changes, ","), " WHERE id=", fmt.Sprint(ud.Id)}

	db := <-DBCH
	defer func(){DBCH <- db}()
	_, _, err := db.Query(strings.Join(stmt, " "))
	if err != nil {
		log.Println(err.Error())
	}
}
