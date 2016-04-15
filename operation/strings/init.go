package strings

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/operation"
)

func init() {
	log.WithField("operation-id", "string.reverse").Debug("Initialised")
	operation.Register("string.reverse", NewStringReverse)
}
