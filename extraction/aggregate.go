package extraction

import (
	log "github.com/Sirupsen/logrus"
	"fmt"
)

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
