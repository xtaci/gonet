package city

import . "db"
import . "types"
import "strings"
import "fmt"
//import "log"

func Flush(city *City) {
	fields, values := SQL_dump(city, "id")
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

func LoadCities(user_id int)(cities []City) {
	stmt := "SELECT * from cities where owner_id = '%v'"

	db := <-DBCH
	defer func(){DBCH <- db}()
	rows, res, err := db.Query(stmt, user_id)
	CheckErr(err)

	for _, row := range rows {
		var city City
		SQL_load(&city, &row, res)
		cities = append(cities, city)
	}

	return
}
