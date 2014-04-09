package protos

import (
	"log"
	"runtime"
)

import (
	"db/stats_tbl"
	"misc/packet"
)

func checkErr(err error) {
	if err != nil {
		funcName, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("ERR:%v,[func:%v,file:%v,line:%v]\n", err, runtime.FuncForPC(funcName).Name(), file, line)
		}

		panic("error occured in Stats Protocol Module")
	}
}

func P_set_adds_req(reader *packet.Packet) []byte {
	tbl, err := PKT_SET_ADDS_REQ(reader)
	//log.Println(tbl)
	checkErr(err)
	stats_tbl.SetAdds(tbl.F_key, tbl.F_value, tbl.F_lang)
	return nil
}

func P_set_update_req(reader *packet.Packet) []byte {
	tbl, err := PKT_SET_UPDATE_REQ(reader)
	//log.Println(tbl)
	checkErr(err)
	stats_tbl.SetUpdate(tbl.F_key, tbl.F_value, tbl.F_lang)
	return nil
}
