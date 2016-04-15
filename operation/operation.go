// Operations execute transformations on the value of an
// Aggregate. For Aggregates that consist of mutliple sub-aggregates
// (e.g. types 'aggregate' or 'auto') the Operation is
// executed on the combined value.
package operation

import (
	log "github.com/Sirupsen/logrus"
)

// Interface type that all Operations need to provide.
type Operation interface {
	// Executes the given Operation a a value and returns
	// the result
	Execute(value string) string
}

var operations map[string]Operation

func init() {
	operations = make(map[string]Operation)
}

func Register(id string, op Operation) {
	log.WithField("operation-id", id).Info("Registered operation")
	operations[id] = op
}

func Get(id string) Operation {
	op, ok := operations[id]
	if !ok {
		return nil
	}

	return op
}
