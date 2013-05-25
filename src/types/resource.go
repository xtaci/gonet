package types

import (
	"encoding/json"
)

type Res struct {
	Id        int32
	Resources map[string]int32
	Version   uint32
}

func (r *Res) Get(name string) int32 {
	return r.Resources[name]
}

func (r *Res) Set(name string, value int32) {
	r.Resources[name] = value
}

func (r *Res) JSON() string {
	val, _ := json.Marshal(r)
	return string(val)
}
