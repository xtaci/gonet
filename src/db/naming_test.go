package db

import "testing"

func TestCamelcase(t *testing.T) {
	if CamelCase("foo_bar") != "FooBar" {
		t.Error("failed foo_bar")
	}

	if CamelCase("_foo_bar") != "FooBar" {
		t.Error("failed _foo_bar")
	}
}

func TestUnderScore(t *testing.T) {
	if UnderScore("FooBar") != "foo_bar" {
		t.Error("failed FooBar")
	}
}
