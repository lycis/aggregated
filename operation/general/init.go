package general

import (
	"github.com/lycis/aggregated/operation"
)

func init() {
	operation.Register("equals", newEquals)
}
