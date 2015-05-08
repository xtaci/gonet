package naming

import (
	"regexp"
	"strings"
)

//----------------------------------------------- "FooBar" => "foo_bar"
var regexp1 = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
var regexp2 = regexp.MustCompile(`([a-z])([A-Z])`)

func UnderScore(str string) string {
	ret := regexp1.ReplaceAllString(str, `${1}_${2}`)
	ret = regexp2.ReplaceAllString(ret, `${1}_${2}`)

	return strings.ToLower(strings.Replace(string(ret), "-", "_", -1))
}

//----------------------------------------------- "foo_bar" -> "FooBar"
var regexp3 = regexp.MustCompile(`^[a-z]|_[a-z]`)

func CamelCase(str string) string {
	str = strings.ToLower(str)

	ret := regexp3.ReplaceAllFunc([]byte(str), func(match []byte) []byte {
		v := strings.TrimLeft(string(match), "_")
		v = strings.ToUpper(v)
		return []byte(v)
	})

	return string(ret)
}

//------------------------------------------------ FNV-1a 32-bit String Hash
func FNV1a(str string) uint32 {
	FNV_prime := uint32(16777619)
	hash := uint32(2166136261)
	octects := []byte(str)

	for _, v := range octects {
		hash = hash ^ uint32(v)
		hash = hash * FNV_prime
	}

	return hash
}
