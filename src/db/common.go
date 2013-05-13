package db

import "log"
import "runtime"

//----------------------------------------------- error checking and logging for database
func CheckErr(err error) {
	if err != nil {
		funcName, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("ERR:%v,[func:%v,file:%v,line:%v]\n", err, runtime.FuncForPC(funcName).Name(), file, line)
		}

		panic("error occured in database module")
	}
}
