package extraction

import (
	log "github.com/Sirupsen/logrus"
)

// extraction for type 'http'
type HttpExtraction struct {
	// URL of the target endpoint
	Url string
}

func (e HttpExtraction) Extract() string {
	log.Debug("http extraction")
	return ""
}