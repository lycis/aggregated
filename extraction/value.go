package extraction

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
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

func (v MultiValue) NestedArray() []interface{} {
	log.WithField("values-cnt", len(v.Values)).Debug("Nesting values")
	resultList := make([]interface{}, len(v.Values))
	for i, subvalue := range v.Values {
		log.WithField("value-num", i).Debug("Processing value")
		if _, ok := subvalue.(SingleValue); ok {
			log.WithField("value-num", i).Debug("single value")
			resultList[i] = subvalue.String()
		} else {
			log.WithField("value-num", i).Debug("multi value")
			resultList[i] = subvalue.(MultiValue).NestedArray()
		}
	}
	return resultList
}
