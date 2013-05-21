package playerdata

import (
	. "db"
	. "types"
)

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	FMT = "%v:%v:%v:%v"
)

//----------------------------------------------- Loading Player Data from DB
func Load(user_id int32) (data *PlayerData,  err error) {
	stmt := "SELECT data FROM buildings where user_id ='%v' LIMIT 1"

	db := <-DBCH
	defer func() { DBCH <- db }()

	rows, _, err := db.Query(stmt, user_id)
	CheckErr(err)

	if len(rows) > 0 {
		err = json.Unmarshal([]byte(rows[0].Str(0)), data)
		CheckErr(err)
	}

	err = errors.New(fmt.Sprint("cannot find building belongs to id:%v", user_id))
	return
}

//----------------------------------------------- Storing Player Data into db
func Store(user_id int32, data *PlayerData) {
	stmt := "UPDATE buildings SET data='%v' WHERE user_id = %v"
	json_estate, err := json.Marshal(*data)
	CheckErr(err)

	db := <-DBCH
	defer func() { DBCH <- db }()
	_, _, err = db.Query(stmt, json_estate, user_id)

	CheckErr(err)
}
