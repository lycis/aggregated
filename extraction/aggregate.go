package extraction

import (
	log "github.com/Sirupsen/logrus"
)

// extraction for type 'aggregate'
type AggregateExtraction struct {
	// Ids of collected sub aggregates
	Ids []string
}

func (e AggregateExtraction) Extract() string {
	log.Debug("aggregate extraction")
	return ""
}