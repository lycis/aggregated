package extraction

import (
	"encoding/json"
	"fmt"
)

type Value interface {
	String() string
}

type SingleValue struct {
	Value string
}

func (v SingleValue) String() string {
	return v.Value
}

type MultiValue struct {
	Values []Value
}

func (v *MultiValue) Add(value Value) {
	v.Values = append(v.Values, value)
}

func (v MultiValue) String() string {
	j, e := json.MarshalIndent(&v, "", "  ")
	if e != nil {
		return fmt.Sprintf("error: %s", e.Error())
	}
	return string(j)
}
