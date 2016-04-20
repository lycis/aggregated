package extraction

import (
	"github.com/lycis/aggregated/extraction"
)

func init() {
	extraction.Register("static", createStaticExtraction)
}

// extraction for type 'static'
type StaticExtraction struct {
	// static value that will be returned
	Value string
}

func (e StaticExtraction) Extract(valueCache map[string]extraction.Value) extraction.Value {
	return extraction.SingleValue{e.Value}
}

func (e StaticExtraction) Dependencies() []string {
	return nil
}

func createStaticExtraction(id string, args interface{}) extraction.Extraction {
	value, ok := args.(string)
	if !ok {
		extraction.PanicParameterError(id, "static")
	}

	extractor := StaticExtraction{
		Value: value,
	}

	return extractor
}
