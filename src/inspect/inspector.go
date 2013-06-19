package inspect

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

import (
	"agent/ipc"
)

func Inspect(id int32, output io.Writer) {
	sess := ipc.QueryOnline(id)
	fmt.Fprintf(output, "%+v\n", sess)
}

func InspectField(id int32, field string, output io.Writer) {
	fields := strings.Split(field, ".")
	fields = fields[1:]

	sess := ipc.QueryOnline(id)
	node := reflect.ValueOf(sess).Elem()

	if sess == nil {
		fmt.Fprintln(output, "user offline")
		return
	}

	for _, v := range fields {
		node = node.FieldByName(v)

		switch node.Kind() {
		case reflect.Ptr, reflect.Interface:
			node = node.Elem()
		}

		if !node.IsValid() {
			fmt.Fprintln(output, "no such field")
			return
		}
	}

	Print(output, node.Interface())
}

func ListAll(output io.Writer) {
	Print(output, ipc.ListAll())
}

func prompt(output io.Writer) {
	fmt.Fprint(output, "gonet> ")
}

func Print(output io.Writer, value interface{}) {
	txt, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		fmt.Fprintln(output, err)
	} else {
		fmt.Fprintln(output, string(txt))
	}
}
