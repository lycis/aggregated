package aggregate

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/configuration"
	"github.com/lycis/aggregated/extraction"
)

var loadedAggregates map[string]*Aggregate

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

// Parse the configuration for Aggregate definitions
// returns the number of loaded aggregates
func LoadAggregates(y configuration.YamlContent) int {
	aggregates := make(map[string]*Aggregate)

	for name, def := range y {
		defer func() {
			if r := recover(); r != nil {
				log.WithError(r.(error)).WithField("id", name).Error("Aggregate definition error")
			}
		}()

		if name != "aggregated" {
			log.WithField("id", name).Debug("Processing aggregate definition")
			aggregate := BuildAggregateFromDefinition(name, def)
			aggregate.Id = name
			log.WithField("id", name).Debug("Definition processed.")

			log.WithField("id", name).Debug("Preparing aggregate")
			aggregate.UpdateExtractor()
			log.WithField("id", name).Debug("Aggregate prepared")

			log.WithFields(log.Fields{
				"id":   name,
				"name": aggregate.Name,
				"type": aggregate.Type,
			}).Info("Discovered aggregate")
			aggregates[name] = &aggregate
		}
	}

	loadedAggregates = aggregates
	return len(loadedAggregates)
}

// Returns the aggregate with the given id if it can be found.
// If not nil will be returned
func GetAggregate(id string) *Aggregate {
	a, ok := loadedAggregates[id]
	if !ok {
		return nil
	}
	return a
}
