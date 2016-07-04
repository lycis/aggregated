package math

import (
	"github.com/lycis/aggregated/extraction"
	"github.com/lycis/aggregated/operation"
	"strconv"
	"fmt"
)

// operation: math.average
// Calculates the average of all subordinate values.
// Will result in an error if a value ist not numeric.
type MathAverage struct{}

func (m MathAverage) Execute(in extraction.Value) extraction.Value {
	
	// no average for single values
	_, ok := in.(extraction.SingleValue)
	if ok {
		return in
	}
	
	mv := in.(extraction.MultiValue)
	
	if len(mv.Values) < 1 {
		return extraction.SingleValue{"0"}
	}
	
	var sum float64
	for _ , v := range mv.Values {
		i, err := strconv.ParseFloat(v.String(), 64)
		if err != nil {
			panic(err)
		}
		
		sum += i 
	}
	
	rVal := extraction.SingleValue{fmt.Sprintf("%f", float64(sum)/float64(len(mv.Values)))};
	return rVal;
}

func newMathAverage(args interface{}) (operation.Operation) {
	return MathAverage{}
}
