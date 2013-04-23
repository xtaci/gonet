package cfg

import "testing"

func TestReadConfig(t *testing.T) {
	m := read_config("config.ini")

	if m["service"] == "" {
		t.Error("load config file failed")
	}
}
