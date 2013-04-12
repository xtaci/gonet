package db

import "testing"

func TestCamelcase(t *testing.T) {
	if camelcase("foo_bar") == "FooBar" && camelcase("_foo_bar") == "FooBar" {
		t.Log("passed")
	} else {
		t.Error("failed")
	}
}
