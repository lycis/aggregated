package configuration

import (
	log "github.com/Sirupsen/logrus"
)

// Used for erros in the service definition
type DefinitionError struct {
	Message string
}

func (e DefinitionError) Error() string {
	return e.Message
}

// General type for data that is retrieved from the
// YAML configuration
type YamlContent map[string]interface{}

// Gives the ServiceDefinition that is represented in this YAML content.
// May also return nil if there is no service definition present.
func (y YamlContent) ServiceDefintion() *ServiceDefinition {
	if _, ok := y["aggregated"]; !ok {
		return nil
	}

	var def ServiceDefinition
	m := y["aggregated"].(map[interface{}]interface{})
	bind, exists := m["bind"].(string)
	if exists {
		def.Bind = bind
	} else {
		def.Bind = ""
	}

	logConfig, ok := m["log"].(map[interface{}]interface{})
	if ok {
		level, exists := logConfig["level"].(string)
		if exists {
			def.Log.Level = level
		} else {
			def.Log.Level = ""
		}
		setLogLevel(level)
	}

	return &def
}

// Adapt the log level
func setLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.WithField("level", log.GetLevel()).Info("log level set")
}
