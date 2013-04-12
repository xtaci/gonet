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

	test.T = time.Date(2013, time.April, 16, 10, 22, 0, 0, time.UTC)
	test.D = "test"
	_, values := sql_dump(&test)

	if values[0] != "'2013-04-16 10:22:00'" {
		t.Error("time dump failed", values[0])
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
}
