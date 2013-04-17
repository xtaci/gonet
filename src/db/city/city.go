package city

import . "db"
import . "types"
import "strings"
import "log"

func Flush(city *City) {
	fields, values := SQL_dump(city)
	stmt := []string{"REPLACE INTO cities(", strings.Join(fields, ","),
		") VALUES (", strings.Join(values, ","), ")"}

	db := <-DBCH
	defer func(){DBCH <- db}()
	_, _, err := db.Query(strings.Join(stmt, " "))
	if err != nil {
		log.Println(err.Error())
	}
}

func New(city *City) {
}
