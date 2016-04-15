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

func (e AutoExtraction) Extract(valueCache map[string]string) string {
	var value string
	for id, v := range valueCache {
		if strings.HasPrefix(id, e.Id) {
			log.WithField("dependency-id", id).Debug("Using dependency in auto extraction")
			value += ", " + v
		}
	}
	log.Debugf("auto := %s", value)
	return value
}

func (e AutoExtraction) Dependencies() []string {
	return nil
}

func createAutoExtraction(id string, args interface{}) extraction.Extraction {
	return AutoExtraction{id}
}