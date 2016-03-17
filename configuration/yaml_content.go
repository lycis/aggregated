package configuration

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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

// Loads the configuration from a YAML file
func FromFile(configFile string) YamlContent {
	log.WithField("file", configFile).Info("Parsing configuration")
	y := loadYamlIntoMap(configFile)
	return y
}

// Loads the configuration from a directory that contains configuration
// files.
// This is done by merging all files into a single temporary config file and
// loading it via LoadFile
func FromDirectory(configDir string) YamlContent {
	tempFile := mergeConfigFiles(configDir)
	y := FromFile(tempFile)
	if err := os.Remove(tempFile); err != nil {
		log.WithError(err).Warn("Could not delete temporary configuration merge file")
	}
	return y
}

// Merges multiple configuration files from a single directory into
// one big config file
func mergeConfigFiles(configDir string) string {
	defer func() {
		if r := recover(); r != nil {
			log.WithError(r.(error)).Error("Merge of config file failed")
			panic(r)
		}
	}()

	files, err := ioutil.ReadDir(configDir)
	if err != nil {
		panic(err)
	}

	// merge all config into one file
	var overallContent []byte
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", configDir, file.Name())
		log.WithField("file", filePath).Info("Mergeing configuraion from file")
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		overallContent = append(overallContent, content...)
	}

	tempFile, err := ioutil.TempFile(".", "aggregated_tmp_config")
	defer tempFile.Close()
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(tempFile.Name(), overallContent, 0666); err != nil {
		panic(err)
	}

	return tempFile.Name()
}

// Loads the content of a YAML file into the internal
// representation as YamlContent
func loadYamlIntoMap(file string) YamlContent {
	m := make(YamlContent)

	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(content, &m)
	if err != nil {
		panic(err)
	}

	return m
}
