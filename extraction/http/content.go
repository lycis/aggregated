package http

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/extraction"
	"io/ioutil"
)

type HttpContentExtraction struct {
	Connection ConnectionData
}

func (e HttpContentExtraction) Extract(valueCache map[string]extraction.Value) extraction.Value {
	log.Debug("http.content extraction")

	r := e.Connection.eval()
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	return extraction.SingleValue{string(body)}
}

func (e HttpContentExtraction) Dependencies() []string {
	return nil
}

func createHttpContentExtraction(id string, args interface{}) extraction.Extraction {
	cd, ok := getConnectionDataFromArgs(args)
	if !ok {
		extraction.PanicParameterError(id, "http")
	}

	extractor := HttpContentExtraction{
		Connection: cd,
	}

	return extractor
}
