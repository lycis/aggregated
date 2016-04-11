package extraction

import (
)

// extraction for type 'static'
type StaticExtraction struct {
	// static value that will be returned
	Value string
}

func (e StaticExtraction) Extract(valueCache map[string]string) string {
	return e.Value
}

func (e StaticExtraction) Dependencies() []string {
	return nil
}

func createStaticExtraction(id string, args interface{}) Extraction {
	value, ok := args.(string)
	if !ok {
		panicParameterError(id, "static")
	}

	extractor := StaticExtraction{
		Value: value,
	}
	
	return extractor
}
