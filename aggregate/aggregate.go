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
	Id          string
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
		aggregate.applyAggregateExtractor()
	case "auto":
		aggregate.applyAutoExtractor()
	default:
		panic(&AggregateDefinitionError{"unsupported type"})
	}
}

// applies an HTTP Extractor to the given aggregate
func (aggregate *Aggregate) applyHttpExtractor() {
	parameters, ok := aggregate.Args.(map[interface{}]interface{})
	if !ok {
		panicParameterError(aggregate.Id, "http")
	}

	extractor := extraction.HttpExtraction{
		Url: parameters["url"].(string),
	}

	aggregate.Extractor = extractor
}

func (a *Aggregate) applyAggregateExtractor() {
	parameters, ok := a.Args.([]interface{})
	if !ok {
		panicParameterError(a.Id, "aggregate")
	}

	var extractor extraction.AggregateExtraction
	for _, sub := range parameters {
		typedSub, ok := sub.(string)
		if !ok {
			panicParameterError(a.Id, "aggregate")
		}
		extractor.Ids = append(extractor.Ids, typedSub)
	}

	a.Extractor = extractor
}

func (a *Aggregate) applyAutoExtractor() {
	extractor := extraction.AutoExtraction{a.Id}
	a.Extractor = extractor
}
