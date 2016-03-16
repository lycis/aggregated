package aggregate

import (
	"github.com/lycis/aggregated/extraction"
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
	Name        string
	Type        string
	Args        interface{}
	OperationId string
	Extractor   extraction.Extraction
}

// Sets Extractor function
func (aggregate *Aggregate) UpdateExtractor() {
	switch aggregate.Type {
	case "http":
		aggregate.applyHttpExtractor()
	case "aggregate":
		//aggregate.applyAggregateExtractor()
	default:
		panic(&AggregateDefinitionError{"unsupported type"})
	}
}

// applies an HTTP Extractor to the given aggregate
func (aggregate *Aggregate) applyHttpExtractor() {
	parameters, ok := aggregate.Args.(map[interface{}]interface{})
	if !ok {
		panic(&AggregateDefinitionError{"invalid definition of parameters for type 'http'"})
	}

	extractor := extraction.HttpExtraction{
		Url: parameters["url"].(string),
	}

	aggregate.Extractor = extractor
}
