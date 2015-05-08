package naming

import (
	"fmt"
	"testing"
)

func TestNaming(t *testing.T) {
	if CamelCase("foo_bar") != "FooBar" {
		t.Error("failed foo_bar")
	}

	if CamelCase("_foo_bar") != "FooBar" {
		t.Error("failed _foo_bar")
	}

	if UnderScore("FooBar") != "foo_bar" {
		t.Error("failed FooBar")
	}

	fmt.Println("FNV-1a Hashing:")
	fmt.Printf("%x\n", FNV1a(""))
	fmt.Printf("%x\n", FNV1a("a"))
	fmt.Printf("%x\n", FNV1a("foobar"))

	if FNV1a("foobar") != 0xbf9cf968 {
		t.Error("FNV hash failed")
	}
}

func BenchmarkFNV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FNV1a("foobar")
	}
}
