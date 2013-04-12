package db

import "time"
import "testing"


type Test struct {
	T time.Time
	A int
	B uint
	C int32
	D string
}

func TestSqlDump(t *testing.T) {
	var test Test

	ti := time.Date(2013, time.April, 16, 10, 22, 0, 0, time.UTC)
	test.T = ti
	test.D = "test"
	_, values := sql_dump(&test)

	if values[0] != "'2013-04-16 10:22:00'" {
		t.Error("time dump failed")
	}

	if values[1] != "'0'" {
		t.Error("int dump failed")
	}

	if values[2] != "'0'" {
		t.Error("uint dump failed")
	}

	if values[3] != "'0'" {
		t.Error("uint dump failed")
	}

	if values[4] != "'test'" {
		t.Error("string dump failed")
	}

	ti2,_ := time.Parse("2006-01-02 15:04:05", "2013-04-16 10:22:00")
	gob,_ := ti2.GobEncode()
	var ti3 time.Time
	ti3.GobDecode(gob)

	if (ti3.Format("2006-01-02 15:04:05") != "2013-04-16 10:22:00") {
		t.Error("time gob failed")
	}
}
