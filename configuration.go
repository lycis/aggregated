package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// This contains the definition of the aggregated service.
// It contains the essentials like which port to bind to and
// the used log level.
type ServiceDefinition struct {
	Bind string
	Log  struct {
		Level string
	}
}

// Merge two ServiceDefinitions. Used if there are multiple "aggregated"
// entries in the config files. The data of subsequent findings is added
// to the first found only if previous definitions did not set that
// inormation already.
func (this *ServiceDefinition) Merge(other ServiceDefinition) {
	if this.Bind == "" {
		this.Bind = other.Bind
	}

	if this.Log.Level == "" {
		this.Log.Level = other.Log.Level
	}
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

// Loads the configuration from a YAML file
func (s ServiceDefinition) FromFile(configFile string) (ServiceDefinition, map[string]Aggregate) {
	log.WithField("file", configFile).Info("Parsing configuration")
	y := loadYamlIntoMap(configFile)

	// parse service configuration
	var sd ServiceDefinition
	if aggregated := y.ServiceDefintion(); aggregated != nil {
		sd = *aggregated
	}

	// parse aggregates
	aggregates := parseAggregates(y)

	return sd, aggregates
}

// Parse the configuration for Aggregate definitions
func parseAggregates(y YamlContent) map[string]Aggregate {
	aggregates := make(map[string]Aggregate)
	
	for name, def := range y {
		defer func() {
			if r := recover(); r != nil {
				log.WithError(r.(error)).Error("Aggregate definition error")
			}	
		}()
		
		if name != "aggregated" {
			aggregates[name] = buildAggregateFromDefinition(name, def)
			log.WithFields(log.Fields {
					"id": name,
					"name": aggregates[name].Name,
					"type": aggregates[name].Type,
			}).Info("Discovered aggregate")
		}
	}
	return aggregates
}

// Loads the configuration from a directory that contains configuration
// files.
// This is done by merging all files into a single temporary config file and
// loading it via LoadFile
func (c ServiceDefinition) FromDirectory(configDir string) (ServiceDefinition, map[string]Aggregate) {
	tempFile := mergeConfigFiles(configDir)
	sd, aggregates := c.FromFile(tempFile)
	if err := os.Remove(tempFile); err != nil {
		log.WithError(err).Warn("Could not delete temporary configuration merge file")
	}
	return sd, aggregates
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
		log.WithField("file", filePath).Info("Loading configuraion from file")
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		overallContent = append(overallContent, content...)
	}

	tempFile, err := ioutil.TempFile(".", "aggregated_tmp_config")
	if err != nil {
		panic(err)
	}
	
	if err := ioutil.WriteFile(tempFile.Name(), overallContent, 0666); err != nil {
		panic(err)
	}
	
	if err := tempFile.Close(); err != nil {
		panic(err)
	}
	
	return tempFile.Name()
}
