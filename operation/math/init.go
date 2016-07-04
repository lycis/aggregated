package math

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/operation"
)

func init() {
	log.WithField("operation-id", "math.average").Debug("Initialised")
	operation.Register("math.average", newMathAverage)
}
