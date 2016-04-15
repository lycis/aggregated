package configuration

// This contains the definition of the aggregated service.
// It contains the essentials like which port to bind to and
// the used log level.
type ServiceDefinition struct {
	Bind string
	Log  struct {
		Level string
	}
}
