package types

import (
	"encoding/json"
)

type Res struct {
	Id        int32
	Resources map[string]int32
}

func (r *Res) Get(name string) {
	return r.Resources[name]
}

func (r *Res) JSON() string {
	val, _ := json.Marshal(r)
	return string(val)
}
