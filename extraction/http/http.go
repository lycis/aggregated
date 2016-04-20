package http

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/extraction"
)

func init() {
	extraction.Register("http", createHttpExtraction)
}

// extraction for type 'http'
type HttpExtraction struct {
	// URL of the target endpoint
	Url string
}

func (e HttpExtraction) Extract(valueCache map[string]extraction.Value) extraction.Value {
	log.Debug("http extraction")
	return extraction.SingleValue{"http"}
}

func (e HttpExtraction) Dependencies() []string {
	return nil
}

func createHttpExtraction(id string, args interface{}) extraction.Extraction {
	parameters, ok := args.(map[interface{}]interface{})
	if !ok {
		extraction.PanicParameterError(id, "http")
	}

	extractor := HttpExtraction{
		Url: parameters["url"].(string),
	}

	return extractor
}
