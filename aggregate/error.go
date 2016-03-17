package aggregate

import (
	"fmt"
)

// Used for errors in the Aggregate definiton
type AggregateDefinitionError struct {
	Message string
}

func (e AggregateDefinitionError) Error() string {
	return e.Message
}

// Panics with an AggregateError. Helper function.
func panicAggregateError(id, message string, inserts ...interface{}) {
	panic(&AggregateDefinitionError{
		Message: fmt.Sprintf(id+": "+message, inserts),
	})
}

func panicParameterError(id, t string) {
	panic(&AggregateDefinitionError{fmt.Sprintf("%s: invalid definition of parameters for type '%s'", id, t)})
}
