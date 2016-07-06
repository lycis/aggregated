package general

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/aggregated/extraction"
	"github.com/lycis/aggregated/operation"
)

// operation: equals
// checks if a aggregate value equals the given value
type Equals struct {
	Value string
}

func (m Equals) Execute(in extraction.Value) extraction.Value {
	log.WithFields(log.Fields{"expected": m.Value, "actual": in.String()}).Info("Comparing values")
	if in.String() == m.Value {
		return extraction.SingleValue{"true"}
	} else {
		return extraction.SingleValue{"false"}
	}
}

func newEquals(args interface{}) operation.Operation {
	log.WithField("args", args).Debug("Building operation 'equals'")
	return Equals{Value: args.(string)}
}
