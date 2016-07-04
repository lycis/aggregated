package http

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/extraction"
	"fmt"
)

// extraction for type 'http'
type HttpStatusExtraction struct {
	Connection ConnectionData
}

func (e HttpStatusExtraction) Extract(valueCache map[string]extraction.Value) extraction.Value {
	log.Debug("http extraction")
	return extraction.SingleValue{fmt.Sprintf("%d", e.Connection.eval().StatusCode)}
}

func (e HttpStatusExtraction) Dependencies() []string {
	return nil
}

func createHttpStatusExtraction(id string, args interface{}) extraction.Extraction {
	parameters, ok := args.(map[interface{}]interface{})
	if !ok {
		extraction.PanicParameterError(id, "http")
	}
	
	cd := ConnectionData {
		Url: parameters["url"].(string),
	}
	
	if m, ok := parameters["method"]; ok {
		cd.Method = m.(string)
	}
	
	if a, ok := parameters["auth"]; ok {
		cd.Headers = append(cd.Headers, fmt.Sprintf("Authorization: %s", a.(string)))
	}

	extractor := HttpStatusExtraction{
		Connection: cd,
	}

	return extractor
}
