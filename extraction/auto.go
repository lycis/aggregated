package extraction

import (
	log "github.com/Sirupsen/logrus"
)

// extraction for type 'aggregate'
type AutoExtraction struct {
	// id of associated Aggregate
	Id string
}

func (e AutoExtraction) Extract() string {
	log.Debug("auto extraction")
	return ""
}
