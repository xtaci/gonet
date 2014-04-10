package gamedata

import (
	"encoding/csv"
	"os"
	"path"
	"strings"
)

import (
	. "helper"
)

//----------------------------------------------- parse & load a game data file into dictionary
func parse(file *os.File) {
	// csv 读取器
	csv_reader := csv.NewReader(file)
	records, err := csv_reader.ReadAll()
	if err != nil {
		ERR("cannot parse csv file.", file.Name(), err)
		return
	}

	// 是否为空档
	if len(records) == 0 {
		ERR("csv file is empty", file.Name())
		return
	}

	// 处理表名
	fi, err := file.Stat()
	if err != nil {
		ERR("cannot stat the file", file.Name())
		return
	}
	tblname := strings.TrimSuffix(fi.Name(), path.Ext(file.Name()))

	// 记录数据, 第一行为表头，因此从第二行开始
	for line := 1; line < len(records); line++ {
		for field := 1; field < len(records[line]); field++ { // 每条记录的第一个字段作为行索引
			_set(tblname, records[line][0], records[0][field], records[line][field])
		}
	}
}
