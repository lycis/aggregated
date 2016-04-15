package extraction

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/extraction"
)

func init() {
	extraction.Register("aggregate", createAggregateExtraction)
}

// extraction for type 'aggregate'
type AggregateExtraction struct {
	// Ids of collected sub aggregates
	Ids []string
}

func (e AggregateExtraction) Extract(valueCache map[string]string) string {
	var value string
	for _, id := range e.Ids {
		v, ok := valueCache[id]
		if !ok {
			log.WithField("error", fmt.Sprintf("dependency %s was not resolved", id)).Error("Failed to resolve value.")
		} else {
			value = fmt.Sprintf("%s, %s", value, v)
		}
	}
	log.Debug(fmt.Sprintf("aggregate extraction := %s", value))
	return value
}

func (e AggregateExtraction) Dependencies() []string {
	return e.Ids
}

func createAggregateExtraction(id string, args interface{}) extraction.Extraction {
	parameters, ok := args.([]interface{})
	if !ok {
		extraction.PanicParameterError(id, "aggregate")
	}

	var extractor AggregateExtraction
	for _, sub := range parameters {
		typedSub, ok := sub.(string)
		if !ok {
			extraction.PanicParameterError(id, "aggregate")
		}
		extractor.Ids = append(extractor.Ids, typedSub)
	}

	return extractor
}
