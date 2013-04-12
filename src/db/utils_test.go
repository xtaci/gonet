package db

import "testing"

func TestCamelcase(t *testing.T) {
	if camelcase("foo_bar") != "FooBar" {
		t.Error("failed foo_bar")
	}

	if camelcase("_foo_bar") != "FooBar" {
		t.Error("failed _foo_bar")
	}
}

func TestUnderscore(t *testing.T) {
	if underscore("FooBar") != "foo_bar" {
		t.Error("failed FooBar")
	}
}
