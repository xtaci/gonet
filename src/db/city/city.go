package city

import . "db"
import . "types"
import "strings"

func Flush(city *City) {
	fields, values := SQL_dump(city)
	stmt := []string{"REPLACE INTO cities(", strings.Join(fields, ","),
		") VALUES (", strings.Join(values, ","), ")"}

	db := <-DBCH
	defer func(){DBCH <- db}()
	_,_, err := db.Query(strings.Join(stmt, " "))
	CheckErr(err)
}

func Create(city *City) {
	fields, values := SQL_dump(city, "id")
	stmt := []string{"INSERT INTO cities(", strings.Join(fields, ","),
		") VALUES (", strings.Join(values, ","), ")"}

	db := <-DBCH
	defer func(){DBCH <- db}()
	_,res, err := db.Query(strings.Join(stmt, " "))
	CheckErr(err)
	city.Id = res.InsertId()
}
