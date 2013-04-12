package db

import "regexp"
import "strings"

// "FooBar" => "foo_bar"
func underscore(str string) string {
	re := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	ret := re.ReplaceAllString(str, `${1}_${2}`)

	re =  regexp.MustCompile(`([a-z])([A-Z])`)
	ret = re.ReplaceAllString(ret, `${1}_${2}`)

	return strings.ToLower(strings.Replace(string(ret), "-", "_", -1))
}
