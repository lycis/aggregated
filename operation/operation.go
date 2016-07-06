// Operations execute transformations on the value of an
// Aggregate. For Aggregates that consist of mutliple sub-aggregates
// (e.g. types 'aggregate' or 'auto') the Operation is
// executed on the combined value.
package operation

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/extraction"
)

// Interface type that all Operations need to provide.
type Operation interface {
	// Executes the given Operation a a value and returns
	// the result
	Execute(value extraction.Value) extraction.Value
}

// This type describes a function that creates an Operation.
// The passed args are directly mapped from the YAML definition
// of the operation
type OperationConstructor func(args interface{})Operation

var operations map[string]OperationConstructor

func init() {
	operations = make(map[string]OperationConstructor)
}

// Registers an operation. The method passed here as "op" has to 
// construct an Operation in some kind.
func Register(id string, op OperationConstructor) {
	log.WithField("operation-id", id).Info("Registered operation")
	operations[id] = op
}

func Get(id string, args interface{}) Operation {
	op, ok := operations[id]
	if !ok {
		return nil
	}

	return op(args)
}
