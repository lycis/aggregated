package extraction

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/extraction"
	"strings"
)

func init() {
	extraction.Register("auto", createAutoExtraction)
}

// extraction for type 'aggregate'
type AutoExtraction struct {
	// id of associated Aggregate
	Id string
}

func (e AutoExtraction) Extract(valueCache map[string]extraction.Value) extraction.Value {
	var value extraction.MultiValue
	for id, v := range valueCache {
		if strings.HasPrefix(id, e.Id) {
			log.WithField("dependency-id", id).Debug("Using dependency in auto extraction")
			value.Add(v)
		}
	}
	log.Debugf("auto := %s", value.String())
	return value
}

func (e AutoExtraction) Dependencies() []string {
	return nil
}

func createAutoExtraction(id string, args interface{}) extraction.Extraction {
	return AutoExtraction{id}
}
