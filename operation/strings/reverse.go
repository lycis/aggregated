// Operation: string.reverse
//
// Returns the given value in reverted order.
//
// Example:
//    abcd => dcba
package strings

type StringReverse struct {}

func (s StringReverse) Execute(in string) (out string) {
	runes := []rune(in)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	
	out = string(runes)
	return
} 