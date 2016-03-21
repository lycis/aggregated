package extraction

import (
	log "github.com/Sirupsen/logrus"
	"strings"
)

// extraction for type 'aggregate'
type AutoExtraction struct {
	// id of associated Aggregate
	Id string
}

func (e AutoExtraction) Extract(valueCache map[string]string) string {
	var value string
	for id, v := range valueCache {
		if strings.HasPrefix(id, e.Id) {
			log.WithField("dependency-id", id).Debug("Using dependency in auto extraction")
			value += ", " + v 
		}
	}
	return value
}
