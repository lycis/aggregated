package extraction

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func init() {
	registeredExtractions = make(map[string]ExtractionConstructor)
}

// General inerface for value extractions
type Extraction interface {
	// Executes the extraction and calculates its value
	Extract(valueCache map[string]string) string

	// Returns a list of dependencies that are injected
	// into the aggregate by this extraction.
	Dependencies() []string
}

// A extraction constructor takes an interface as argument.
// It has to be casted and checked by the constructor function and
// return a Extraction instance.
// Args are the values defined in the 'args' element of the aggregate
// definition.
type ExtractionConstructor func(id string, args interface{}) Extraction

var registeredExtractions map[string]ExtractionConstructor

// Registers a new Extraction type. Call this in your
// init function.
func Register(name string, constructor ExtractionConstructor) {
	log.WithField("type", name).Info("Registered extraction")
	registeredExtractions[name] = constructor
}

// Returns a constructor for the given Extraction type or nil
// for unknown types
func Get(name string) ExtractionConstructor {
	constructor, ok := registeredExtractions[name]
	if !ok {
		return nil
	}

	return constructor
}

type ParameterError struct {
	Message string
}

func (e ParameterError) Error() string {
	return e.Message
}

// Panics with an error that indicates that a parameter definition
// of an Extraction type was violated
func PanicParameterError(id, t string) {
	panic(&ParameterError{fmt.Sprintf("%s: invalid definition of parameters for type '%s'", id, t)})
}
