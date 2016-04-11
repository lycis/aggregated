package extraction

import (
	log "github.com/Sirupsen/logrus"
)

// extraction for type 'http'
type HttpExtraction struct {
	// URL of the target endpoint
	Url string
}

func (e HttpExtraction) Extract(valueCache map[string]string) string {
	log.Debug("http extraction")
	return ""
}

func (e HttpExtraction) Dependencies() []string {
	return nil
}

func createHttpExtraction(id string, args interface{}) Extraction {
	parameters, ok := args.(map[interface{}]interface{})
	if !ok {
		panicParameterError(id, "http")
	}

	extractor := HttpExtraction{
		Url: parameters["url"].(string),
	}
	
	return extractor
}
