package db

import . "types"
import "strings"
import "log"

func CityFlush(city *City) {
	fields, values := sql_dump(city)
	stmt := []string{"REPLACE INTO cities(", strings.Join(fields, ","),
		") VALUES (", strings.Join(values, ","), ")"}

	db := <-DBCH
	defer func(){DBCH <- db}()
	_, _, err := db.Query(strings.Join(stmt, " "))
	if err != nil {
		log.Println(err.Error())
	}
}

func CityNew(city *City) {
}
