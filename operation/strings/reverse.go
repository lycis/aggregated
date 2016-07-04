package strings

import (
	"github.com/lycis/aggregated/extraction"
	"github.com/lycis/aggregated/operation"
)

// Operation: string.reverse
//
// Returns the given value in reverted order.
//
// Example:
//    abcd => dcba
type StringReverse struct{}

func (s StringReverse) Execute(in extraction.Value) (extraction.Value) {
	return reverseValue(in)
}

func reverseValue(in extraction.Value) (out extraction.Value) {
	sv, ok := in.(extraction.SingleValue)
	if ok {
		out = reverseSingleValue(sv)
	} else {
		mv, ok := in.(extraction.MultiValue)
		if ok {
			out = reverseMultiValue(mv)
		}
	}
	return
}

func reverseSingleValue(in extraction.SingleValue) (extraction.SingleValue) {
	runes := []rune(in.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return extraction.SingleValue{string(runes)}	
}

func reverseMultiValue(in extraction.MultiValue) (extraction.MultiValue) {
	var out extraction.MultiValue
	for _, v := range in.Values {
		out.Add(reverseValue(v))
	}
	return out
}

func NewStringReverse(args interface{}) operation.Operation {
	return StringReverse{}
}
