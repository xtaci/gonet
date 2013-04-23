package db

import "regexp"
import "strings"

// "FooBar" => "foo_bar"
var regexp1 = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
var regexp2 = regexp.MustCompile(`([a-z])([A-Z])`)

func UnderScore(str string) string {
	ret := regexp1.ReplaceAllString(str, `${1}_${2}`)
	ret = regexp2.ReplaceAllString(ret, `${1}_${2}`)

	return strings.ToLower(strings.Replace(string(ret), "-", "_", -1))
}

// "foo_bar" -> "FooBar"
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
