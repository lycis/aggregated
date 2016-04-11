package extraction

// General inerface for value extractions
type Extraction interface {
	Extract(valueCache map[string]string) string
}

// A extraction constructor takes an interface as argument. 
// It has to be casted and checked by the constructor function and
// return a Extraction instance.
// Args are the values defined in the 'args' element of the aggregate
// definition.
type ExtractionConstructor func(args interface{})Extraction

var registeredExtractions map[string]ExtractionConstructor

// Registers a new Extraction type. Call this in your
// init function.
func Register(name string, constructor ExtractionConstructor) {
	registeredExtractions[name] = constructor
}

// Returns a constructor for the given Extraction type or nil
// for unknown types
func Get(name string) ExtractionConstructor {
	constructor, ok := registeredExtractions[name]
	if !ok {
		return nil
	}
	
	return constructor
}
