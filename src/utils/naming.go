package utils

import "regexp"
import "strings"

// "FooBar" => "foo_bar"
func UnderScore(str string) string {
	re := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	ret := re.ReplaceAllString(str, `${1}_${2}`)

	re = regexp.MustCompile(`([a-z])([A-Z])`)
	ret = re.ReplaceAllString(ret, `${1}_${2}`)

	return strings.ToLower(strings.Replace(string(ret), "-", "_", -1))
}

// "foo_bar" -> "FooBar"
func CamelCase(str string) string {
	str = strings.ToLower(str)
	re := regexp.MustCompile(`^[a-z]|_[a-z]`)

	ret := re.ReplaceAllFunc([]byte(str), func(match []byte) []byte {
		v := strings.TrimLeft(string(match), "_")
		v = strings.ToUpper(v)
		return []byte(v)
	})

	return string(ret)
}
