package main

import (
	"fmt"
)

// An Aggregate contains all the information necessary to
// evaluate its value. This element represents one aggregate
// entry in the config file.
//
// Example:
//   my-aggregate:
//     name: My Aggregate
//     type: http
//     args:
//       url: http://example.com/foobar
//
type Aggregate struct {
	Name string
	Type string
	Args interface{}
	Operation string
}

// Builds an aggregate from its definition out of the config file
// In the example the id would be "my aggregate" and the definition
// would be passed as map
func buildAggregateFromDefinition(id string, i interface{}) Aggregate {
	var aggregate Aggregate
	
	def, ok := i.(map[interface{}]interface{})
	if !ok {
		panicAggregateError(id, "wrong definition")
	}
	
	name, ok := def["name"].(string)
	if !ok {
		panicAggregateError(id, "name is not a string")
	}
	aggregate.Name = name
	
	xtype, ok := def["type"].(string)
	if !ok {
		panicAggregateError(id, "type is not a string")
	}
	aggregate.Type = xtype
	
	args, ok := def["args"]
	if ok {
		aggregate.Args = args
	}
	
	untypedOp, ok := def["operation"]
	if ok {
		operation, ok := untypedOp.(string)
		if !ok {
			panicAggregateError(id, "operation is not a string")
		}
		aggregate.Operation = operation
	}
	
	return aggregate
}

// Used for errors in the Aggregate definiton
type AggregateDefinitionError struct {
	Message string
}

func (e AggregateDefinitionError) Error() string {
	return e.Message
}

// Panics with an AggregateError. Helper function.
func panicAggregateError(id, message string, inserts... interface{}) {
	panic(&AggregateDefinitionError{
				Message: fmt.Sprintf(id + ": " + message, inserts),
		})
}