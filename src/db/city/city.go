package city

import . "db"
import . "types"
import "strings"
import "fmt"

func Flush(city *City) {
	fields, values := SQL_dump(city)
	changes := SQL_set_clause(fields,values)

	stmt := []string{"UPDATE cities SET ", strings.Join(changes, ","), " WHERE id=", fmt.Sprint(city.Id)}

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
