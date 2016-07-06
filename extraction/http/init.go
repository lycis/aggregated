package http

import (
	"github.com/lycis/aggregated/extraction"
)

func init() {
	extraction.Register("http.status", createHttpStatusExtraction)
	extraction.Register("http.content", createHttpContentExtraction)
}
