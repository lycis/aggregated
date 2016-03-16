package extraction

// General inerface for value extractions
type Extraction interface {
	Extract() string
}