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
